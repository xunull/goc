package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"github.com/xunull/goc/enhance/stringsx"
	"github.com/xunull/goc/traverse"
	"path/filepath"
	"strconv"
	"strings"
)

type GitDiffItem struct {
	Dir            string
	FileChangedStr string
	InsertStr      string
	DeleteStr      string
	FileChangeNum  int
	InsertNum      int
	DeleteNum      int
}

func (g *GitApi) getDiffItemWithTargetExts(str string, exts map[string]string) *GitDiffItem {
	t := &GitDiffItem{
		Dir: g.Dir,
	}
	insertCount := 0
	deleteCount := 0
	fileCount := 0
	lines := stringsx.SplitLines(str)
	for _, line := range lines {
		items := strings.Fields(line)
		name := items[2]
		insert, err := strconv.Atoi(items[0])
		if err != nil {
			insert = 0
		}
		del, err := strconv.Atoi(items[1])
		if err != nil {
			del = 0
		}

		ext := filepath.Ext(name)
		if _, ok := exts[ext]; ok {
			insertCount += insert
			deleteCount += del
			if insert > 0 || del > 0 {
				fileCount += 1
			}
		}
	}
	t.InsertNum = insertCount
	t.InsertStr = strconv.Itoa(insertCount)
	t.DeleteNum = deleteCount
	t.DeleteStr = strconv.Itoa(deleteCount)
	t.FileChangeNum = fileCount
	t.FileChangedStr = strconv.Itoa(fileCount)

	return t
}

func (g *GitApi) GetGitDiffItem(str string, opts ...Option) *GitDiffItem {
	o := g.getOption(opts...)

	if o.TargetExt != "" {
		m := map[string]string{
			o.TargetExt: o.TargetExt,
		}
		return g.getDiffItemWithTargetExts(str, m)
	} else if o.OnlyBackLanguage {
		m := traverse.CommonBackLanguageExt
		return g.getDiffItemWithTargetExts(str, m)
	} else if o.OnlyFrontLanguage {
		m := traverse.CommonFrontLanguageExt
		return g.getDiffItemWithTargetExts(str, m)
	}

	sl := strings.Split(str, ",")
	if len(sl) == 3 {
		t := &GitDiffItem{
			FileChangedStr: strings.TrimSpace(sl[0]),
			InsertStr:      strings.TrimSpace(sl[1]),
			DeleteStr:      strings.TrimSpace(sl[2]),
		}

		fc, err := strconv.Atoi(strings.Fields(sl[0])[0])
		if err == nil {
			t.FileChangeNum = fc
		}
		in, err := strconv.Atoi(strings.Fields(sl[1])[0])
		if err == nil {
			t.InsertNum = in
		}
		dn, err := strconv.Atoi(strings.Fields(sl[2])[0])
		if err == nil {
			t.DeleteNum = dn
		}
		return t
	} else if len(sl) == 2 {
		t := &GitDiffItem{
			FileChangedStr: strings.TrimSpace(sl[0]),
		}
		if strings.Contains(sl[1], "insertions") {
			t.InsertStr = strings.TrimSpace(sl[1])
			in, err := strconv.Atoi(strings.Fields(sl[1])[0])
			if err == nil {
				t.InsertNum = in
			}
		} else {
			t.DeleteStr = strings.TrimSpace(sl[1])
			dn, err := strconv.Atoi(strings.Fields(sl[1])[0])
			if err == nil {
				t.DeleteNum = dn
			}
		}
		return t
	} else {
		return &GitDiffItem{}
	}
}

func (g *GitApi) GetDiffCommitAndWorkTree(commit string, ops ...Option) (string, error) {
	o := g.getOption(ops...)
	cmd := []string{"git", "diff", commit}
	if o.ShortStat {
		cmd = append(cmd, "--shortstat")
	} else {
		cmd = append(cmd, "--stat")
	}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))

	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) GetDiffTwoCommit(start, end string, opts ...Option) (string, error) {
	o := g.getOption(opts...)
	cmd := []string{"git", "diff", start, end}
	if o.ShortStat {
		cmd = append(cmd, "--shortstat")
	} else {
		cmd = append(cmd, "--numstat")
	}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))

	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}
