package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"strings"
)

func (g *GitApi) BareClone(remote string) (string, error) {

	cmd := []string{"git", "clone", "--bare", remote, g.Dir}
	cmdRes := commandx.RunCommand(cmd)
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}
