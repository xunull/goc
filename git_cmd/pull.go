package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"strings"
)

func (g *GitApi) GitPull(opts ...Option) (string, error) {
	_ = g.getOption(opts...)
	cmd := []string{"git", "pull"}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String() + "\n" + cmdRes.Stderr.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}
