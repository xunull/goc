package git_cmd

import (
	"bufio"
	"github.com/xunull/goc/commandx"
	"io"
	"strings"
)

func (g *GitApi) WorktreeIsClean() (bool, error) {
	cmd := []string{"git", "status"}

	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		if strings.Contains(cmdRes.Stdout.String(), "working tree clean") {
			return true, nil
		} else if strings.Contains(cmdRes.Stdout.String(), "No commits yet") {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, cmdRes.Err
	}
}

func (g *GitApi) CheckUpToDate(all bool) (bool, *commandx.CommandResult) {
	var cmdOne []string
	if all {
		cmdOne = []string{"git", "fetch", "--all"}
	} else {
		cmdOne = []string{"git", "fetch"}
	}

	cmdTwo := []string{"git", "status"}
	success, res := commandx.RunCommandForLast([][]string{cmdOne, cmdTwo}, commandx.WithDir(g.Dir))
	if success {
		if strings.Contains(res.Stdout.String(), "up to date") {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return false, res
	}
}

func (g *GitApi) StatusPorcelain() ([]string, error) {
	cmd := []string{"git", "status", "--porcelain=1"}

	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		lines := make([]string, 0)
		buff := bufio.NewReader(&cmdRes.Stdout)
		for {
			line, err := buff.ReadString('\n')
			if io.EOF == err {
				break
			}
			if err != nil {
				return nil, err
			}
			lines = append(lines, line)
		}
		return lines, nil
	} else {
		return nil, cmdRes.Err
	}
}

func (g *GitApi) CheckForgetPush() (bool, error) {
	cmd := []string{"git", "status"}

	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		if strings.Contains(cmdRes.Stdout.String(), "is ahead of") {
			return true, nil
		} else {

			return false, nil
		}
	} else {
		return false, cmdRes.Err
	}
}
