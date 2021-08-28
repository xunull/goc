package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"strings"
)

func (g *GitApi) IsWorktreeClean() (bool, error) {
	cmd := []string{"git", "status"}

	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.Contains(cmdRes.Stdout.String(), "working tree clean"), nil
	} else {
		return false, cmdRes.Err
	}
}
