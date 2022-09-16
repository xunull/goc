package git_cmd

import (
	"github.com/xunull/goc/commandx"
)

func (g *GitApi) FetchRemote(remote string) (bool, *commandx.CommandResult) {
	cmd := []string{"git", "fetch", remote}
	res := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	return res.Success, res
}

func (g *GitApi) FetchAll() (bool, *commandx.CommandResult) {
	cmd := []string{"git", "fetch", "--all"}
	res := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	return res.Success, res
}
