package git_cmd

import (
	"github.com/xunull/goc/commandx"
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
