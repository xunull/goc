package git_cmd

import (
	"fmt"
	"github.com/xunull/goc/commandx"
	"strings"
)

func (g *GitApi) CreateBranch(name string, ops ...Option) (string, error) {
	_ = g.getOption(ops...)
	cmd := []string{"git", "branch", name}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) CheckoutRemoteBranch(local, remote, branch string, ops ...Option) (string, error) {
	_ = g.getOption(ops...)
	cmd := []string{"git", "checkout", "-b", local, fmt.Sprintf("%s/%s", remote, branch)}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) PushU(local, remote string, ops ...Option) (string, error) {
	_ = g.getOption(ops...)
	cmd := []string{"git", "push", "-u", remote, local}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) SetUpstreamTo(remote, remoteBranch, localBranch string, ops ...Option) (string, error) {
	_ = g.getOption(ops...)

	t := fmt.Sprintf("--set-upstream-to=%s/%s", remote, remoteBranch)

	cmd := []string{"git", "branch", t, localBranch}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}
