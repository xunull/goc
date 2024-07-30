package commandx

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"os/exec"
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

func RunCommand(command []string, ops ...Option) *CommandResult {
	d := &option{}
	for _, o := range ops {
		o(d)
	}

	if d.Timeout == 0 {
		d.Timeout = time.Second * 20
	}

	ctx, cancel := context.WithTimeout(context.Background(), d.Timeout)
	defer cancel()
	var cmd *exec.Cmd
	if len(command) > 1 {
		cmd = exec.CommandContext(ctx, command[0], command[1:]...)
	} else {
		cmd = exec.CommandContext(ctx, command[0])
	}

	cmd.Env = os.Environ()

	if d.Dir != "" {
		cmd.Dir = d.Dir
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err == nil {
		return &CommandResult{
			Status:  0,
			Stdout:  stdout,
			Stderr:  stderr,
			Success: true,
		}
	} else {
		r := &CommandResult{
			Stdout:  stdout,
			Stderr:  stderr,
			Success: false,
			Err:     err,
		}

		var ex *exec.ExitError
		if errors.As(err, &ex) {
			r.Status = ex.ExitCode()
		}
		return r
	}
}
