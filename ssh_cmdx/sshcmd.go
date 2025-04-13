// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ssh_cmdx

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

// Ping 检查是否能够连接到远程机器
func (c *Client) Ping(host string) error {
	client, _, err := c.Connect(host)
	if err != nil {
		return fmt.Errorf("failed to connect %s: %v", host, err)
	}
	return client.Close()
}

func (c *Client) wrapCommands(cmds ...string) string {
	// 多个命令拼成一行
	cmdJoined := strings.Join(cmds, "; ")
	// 如果不需要sudo 或者用户已经是root 那么就直接返回命令
	if !c.Option.sudo || c.Option.user == defaultUsername {
		return cmdJoined
	}

	// Escape single quotes in cmd, fix https://github.com/labring/sealos/issues/4424
	// e.g. echo 'hello world' -> `sudo -E /bin/bash -c 'echo "hello world"'`
	// 将单引号替换为双引号
	// 然后在加上sudo执行改命令
	cmdEscaped := strings.ReplaceAll(cmdJoined, `'`, `"`)
	return fmt.Sprintf("sudo -E /bin/bash -c '%s'", cmdEscaped)
}

// CmdAsyncWithContext 异步执行cmd
func (c *Client) CmdAsyncWithContext(ctx context.Context, host string, cmds ...string) error {
	cmd := c.wrapCommands(cmds...)
	Debug("start to exec `%s` on %s", cmd, host)
	// 连接到远端机器
	client, session, err := c.Connect(host)
	if err != nil {
		return fmt.Errorf("connect error: %v", err)
	}
	// 关闭session 关闭client
	defer client.Close()
	defer session.Close()
	stdout, err := session.StdoutPipe()
	if err != nil {
		return fmt.Errorf("stdout pipe %s: %v", host, err)
	}
	stderr, err := session.StderrPipe()
	if err != nil {
		return fmt.Errorf("stderr pipe %s: %v", host, err)
	}
	stdin, err := session.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe %s: %v", host, err)
	}
	out := autoAnswerWriter{
		in:        stdin,
		answer:    []byte(c.password + "\n"),
		condition: isSudoPrompt,
	}
	eg, _ := errgroup.WithContext(context.Background())
	// 处理stdout
	eg.Go(func() error { return c.handlePipe(host, stderr, &out, c.stdout) })
	// 处理stderr
	eg.Go(func() error { return c.handlePipe(host, stdout, &out, c.stdout) })

	errCh := make(chan error, 1)
	go func() {
		errCh <- func() error {
			// 执行命令
			if err := session.Start(cmd); err != nil {
				return fmt.Errorf("start command `%s` on %s: %v", cmd, host, err)
			}
			// 等上面处理stdout和stderr完毕
			if err = eg.Wait(); err != nil {
				return err
			}
			// 等待session完毕
			if err = session.Wait(); err != nil {
				return fmt.Errorf("run command `%s` on %s, output: %s, error: %v,", cmd, host, out.b.String(), err)
			}
			return nil
		}()
	}()

	select {
	case <-ctx.Done():
		// 被外界结束
		return ctx.Err()
	case err = <-errCh:
		// 上面的go方法执行完毕
		return err
	}

}

// CmdAsync not actually asynchronously, just print output asynchronously
func (c *Client) CmdAsync(host string, cmds ...string) error {
	ctx, cancel := GetTimeoutContext()
	defer cancel()
	return c.CmdAsyncWithContext(ctx, host, cmds...)
}

// Cmd 同步的执行一个命令
func (c *Client) Cmd(host, cmd string) ([]byte, error) {
	cmd = c.wrapCommands(cmd)
	Debug("start to exec `%s` on %s", cmd, host)
	client, session, err := c.Connect(host)
	if err != nil {
		return nil, fmt.Errorf("failed to create ssh session for %s: %v", host, err)
	}
	defer client.Close()
	defer session.Close()
	in, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	b := autoAnswerWriter{
		in:        in,
		answer:    []byte(c.password + "\n"),
		condition: isSudoPrompt,
	}
	session.Stdout = &b
	session.Stderr = &b
	// Run里面包含了Start和Wait
	err = session.Run(cmd)
	return b.b.Bytes(), err
}

type withPrefixWriter struct {
	prefix  string
	newline bool
	w       io.Writer
	mu      sync.Mutex
}

func (w *withPrefixWriter) Write(p []byte) (int, error) {
	p = append([]byte(w.prefix), p...)
	if w.newline {
		if p[len(p)-1] != byte('\n') {
			p = append(p, '\n')
		}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.w.Write(p)
}

type autoAnswerWriter struct {
	b          bytes.Buffer
	in         io.Writer
	showPrompt bool
	answer     []byte
	condition  func([]byte) bool
	mu         sync.Mutex
}

func (w *autoAnswerWriter) Write(p []byte) (int, error) {
	if w.in != nil && w.condition != nil && w.condition(p) {
		_, err := w.in.Write(w.answer)
		if err != nil {
			return 0, err
		}
		if !w.showPrompt {
			return len(p), nil
		}
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.b.Write(p)
}

func isSudoPrompt(p []byte) bool {
	return bytes.HasPrefix(p, []byte("[sudo] password for ")) && bytes.HasSuffix(p, []byte(": "))
}

// handlePipe 处理stdout stderr的输出
func (c *Client) handlePipe(host string, pipe io.Reader, out io.Writer, isStdout bool) error {
	r := bufio.NewReader(pipe)
	writers := []io.Writer{out}
	if isStdout {
		writers = append(writers, &withPrefixWriter{prefix: host + "\t", newline: true, w: os.Stdout})
	}
	w := io.MultiWriter(writers...)
	var line []byte
	for {
		b, err := r.ReadByte()
		if err != nil {
			// 输出完毕
			if err == io.EOF {
				return nil
			}
			return err
		}
		line = append(line, b)
		// 输出完一行
		if b == byte('\n') || isSudoPrompt(line) {
			// ignore any writer error
			_, _ = w.Write(line)
			line = make([]byte, 0)
			continue
		}
	}
}
