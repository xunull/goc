package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"os"
	"path"
	"strings"
)

func IsGitRepo() bool {
	cmd := []string{"git", "rev-parse", "--is-inside-work-tree"}
	res := commandx.RunCommand(cmd)
	if res.Success {
		return strings.TrimSpace(res.Stdout.String()) == "true"
	} else {
		return false
	}
}

func IsTargetGitRepo(target string) bool {
	cmd := []string{"git", "rev-parse", "--is-inside-work-tree"}
	res := commandx.RunCommand(cmd, commandx.WithDir(target))
	if res.Success {
		return strings.TrimSpace(res.Stdout.String()) == "true"
	} else {
		return false
	}
}

func GetWorkDir() (string, error) {
	cmd := []string{"git", "rev-parse", "--show-toplevel"}
	res := commandx.RunCommand(cmd)
	if res.Success {
		return strings.TrimSpace(res.Stdout.String()), nil
	} else {
		return "", res.Err
	}
}

func GetGitDir() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	cmd := []string{"git", "rev-parse", "--git-dir"}
	res := commandx.RunCommand(cmd)
	if res.Success {
		return path.Join(pwd, strings.TrimSpace(res.Stdout.String())), nil
	} else {
		return "", res.Err
	}
}

func IsInGitDir() bool {
	res := commandx.RunCommand([]string{"git", "rev-parse", "--is-inside-git-dir"})
	if res.Success {
		return strings.TrimSpace(res.Stdout.String()) == "true"
	} else {
		return false
	}
}

func IsBareRepo() bool {
	res := commandx.RunCommand([]string{"git", "rev-parse", "--is-bare-repository"})
	if res.Success {
		return strings.TrimSpace(res.Stdout.String()) == "true"
	} else {
		return false
	}
}
