package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"strings"
)

func (g *GitApi) GitAdd(target []string, opts ...Option) (string, error) {
	o := g.getOption(opts...)
	cmd := []string{"git", "add"}
	if o.AddAll {
		cmd = append(cmd, "-A")
	} else {
		cmd = append(cmd, target...)
	}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}
