package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"github.com/xunull/goc/commonx"
	"regexp"
	"strings"
)

type GitRemoteItem struct {
	ItemBase
	Name string
	Url  string
}

type RepoRemotes struct {
	ItemBase
	Remotes []GitRemoteItem
}

func (r *RepoRemotes) IsHaveRemote(name string) bool {
	for _, item := range r.Remotes {
		if item.Name == name {
			return true
		}
	}
	return false
}

func (g *GitApi) AddRemote(name, url string, ops ...Option) (string, error) {
	_ = g.getOption(ops...)
	cmd := []string{"git", "remote", "add", name, url}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) PushMirror(name string, ops ...Option) (string, error) {
	_ = g.getOption(ops...)
	cmd := []string{"git", "push", name, "--mirror"}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) GetRepoRemotes() (*RepoRemotes, error) {

	if str, err := g.getRemotes(); err == nil {

		if str == "" {
			// no remote
			return nil, nil
		}

		m := make(map[string]string)
		lines := strings.Split(str, "\n")
		for _, line := range lines {
			list := strings.Fields(line)
			m[list[0]] = list[1]
		}
		res := make([]GitRemoteItem, 0, len(m))
		for k, v := range m {
			t := GitRemoteItem{
				Name: k,
				Url:  v,
			}
			t.Dir = g.Dir
			res = append(res, t)
		}
		return &RepoRemotes{
			Remotes: res,
			ItemBase: ItemBase{
				Dir: g.Dir,
			},
		}, nil
	} else {
		return nil, err
	}

}

func (g *GitApi) getRemotes(ops ...Option) (string, error) {
	_ = g.getOption(ops...)
	cmd := []string{"git", "remote", "-v"}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

var usernameReg, _ = regexp.Compile(`^.*/(?P<username>.*)/.*$`)

func (g *GitApi) GetRemoteUsername() string {
	rr, err := g.GetRepoRemotes()
	commonx.CheckErrOrFatal(err)

	if rr != nil && (rr.IsHaveRemote("origin")) {
		for _, item := range rr.Remotes {

			if item.Name == "origin" {
				temp := usernameReg.FindStringSubmatch(item.Url)
				if len(temp) == 2 {
					return temp[1]
				}

			}

		}
	}

	return ""
}
