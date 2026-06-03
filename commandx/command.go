package commandx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"mvdan.cc/sh/v3/interp"
	"mvdan.cc/sh/v3/syntax"
	"os"
	"os/exec"
	"strings"
	"time"
)

type CommandResult struct {
	Stdout  bytes.Buffer
	Stderr  bytes.Buffer
	Status  int
	Success bool
	Err     error
}

func (r *CommandResult) OutputOrFatal() {
	if !r.Success || r.Status != 0 {
		log.Error().Err(r.Err).Msg("")
		fmt.Printf("%s\n", r.Stdout.String())
		fmt.Printf("%s\n", r.Stderr.String())
		os.Exit(1)
	} else {
		fmt.Printf("%s\n", r.Stdout.String())
		fmt.Printf("%s\n", r.Stderr.String())
	}
}

func (r *CommandResult) SimpleLog() string {
	return fmt.Sprintf("code: %d\nstdout: %s\nstderr: %s\n", r.Status, r.Stdout.String(), r.Stderr.String())
}

// RunCommandForLast run commands and return last result
func RunCommandForLast(commands [][]string, ops ...Option) (bool, *CommandResult) {
	f := true
	var res *CommandResult
	for _, command := range commands {
		res = RunCommand(command, ops...)
		if res.Success {
			continue
		} else {
			f = false
			break
		}
	}
	return f, res
}

// HasBash 保留兼容：原语义是检查系统是否有 bash。
// 现在底层改用 mvdan.cc/sh 的内置解释器，已经不依赖系统 bash，
// 但维持函数行为以兼容外部调用方的判断逻辑。
func HasBash() bool {
	_, err := exec.LookPath("bash")
	return err == nil
}

// RunBashCommand 把入参当作一段 shell 源码交给内置解释器执行，
// 支持管道、重定向、变量展开等 bash 语法。
func RunBashCommand(command string, ops ...Option) *CommandResult {
	return runScript(command, ops...)
}

// RunCommand 按 argv 形式执行命令。
// 内部会把参数做 POSIX shell 转义后拼成一行脚本，再交给解释器执行，
// 因此参数里的空格、引号、$ 等不会被二次解释，等价于直接 exec 的语义。
func RunCommand(command []string, ops ...Option) *CommandResult {
	if len(command) == 0 {
		return &CommandResult{
			Success: false,
			Status:  -1,
			Err:     errors.New("commandx: empty command"),
		}
	}
	return runScript(joinShellArgs(command), ops...)
}

func runScript(script string, ops ...Option) *CommandResult {
	d := &option{}
	for _, o := range ops {
		o(d)
	}
	if d.Timeout == 0 {
		d.Timeout = time.Second * 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()

	file, err := syntax.NewParser().Parse(strings.NewReader(script), "")
	if err != nil {
		return &CommandResult{
			Success: false,
			Status:  -1,
			Err:     fmt.Errorf("commandx: parse error: %w", err),
		}
	}

	res := &CommandResult{}
	var stderrW io.Writer = &res.Stderr
	if d.RedirectStderr {
		stderrW = &res.Stdout
	}

	runnerOpts := []interp.RunnerOption{
		interp.StdIO(nil, &res.Stdout, stderrW),
	}
	if d.Dir != "" {
		runnerOpts = append(runnerOpts, interp.Dir(d.Dir))
	}

	runner, err := interp.New(runnerOpts...)
	if err != nil {
		res.Success = false
		res.Status = -1
		res.Err = fmt.Errorf("commandx: new runner: %w", err)
		return res
	}

	runErr := runner.Run(ctx, file)
	if runErr == nil {
		res.Success = true
		res.Status = 0
		return res
	}

	res.Success = false
	res.Err = runErr
	var es interp.ExitStatus
	if errors.As(runErr, &es) {
		res.Status = int(es)
	} else {
		res.Status = -1
	}
	return res
}

// joinShellArgs 把 argv 拼成一行经过 POSIX 转义的 shell 字符串。
func joinShellArgs(args []string) string {
	parts := make([]string, len(args))
	for i, a := range args {
		parts[i] = shellQuote(a)
	}
	return strings.Join(parts, " ")
}

// shellQuote 对单个参数做 POSIX 单引号转义；
// 仅在含有需要转义的字符时才加引号，保持简单参数的可读性。
func shellQuote(s string) string {
	if s == "" {
		return "''"
	}
	for _, r := range s {
		safe := (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '_' || r == '-' || r == '/' || r == '.' ||
			r == ':' || r == '=' || r == '@' || r == '+' ||
			r == ',' || r == '%'
		if !safe {
			return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
		}
	}
	return s
}
