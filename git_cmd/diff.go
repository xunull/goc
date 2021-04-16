package git_cmd

import (
	"github.com/xunull/goc/commandx"
	"github.com/xunull/goc/enhance/stringx"
	"github.com/xunull/goc/lang_ext"
	"path/filepath"
	"sort"
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

type GitDiffItemListStat struct {
	FileCount   int
	InsertCount int
	DeleteCount int
}

type GitDiffShotStat struct {
	FileChange  int
	InsertCount int
	DeleteCount int
}

type GitDiffNumStat struct {
	FileChangedList []string
}

type GitDiffItemList []*GitDiffItem

func (g GitDiffItemList) Sort() {
	sort.Slice(g, func(i, j int) bool {
		if g[i].InsertNum != g[j].InsertNum {
			return g[i].InsertNum > g[j].InsertNum
		} else {
			if g[i].DeleteNum != g[j].DeleteNum {
				return g[i].DeleteNum > g[j].DeleteNum
			} else {
				return g[i].FileChangeNum > g[j].FileChangeNum
			}
		}
	})
}

func (g *GitApi) GetGitDiffNumStat(start, end string, opts ...Option) (*GitDiffNumStat, error) {
	res, err := g.GetDiffTwoCommit(start, end, WithNumStat())
	if err != nil {
		return nil, err
	}

	result := new(GitDiffNumStat)

	lines := strings.Split(res, "\n")
	result.FileChangedList = make([]string, 0, len(lines))
	for _, line := range lines {
		t := strings.Fields(line)
		fp := t[2]
		result.FileChangedList = append(result.FileChangedList, fp)
	}
	return result, nil
}

func (g *GitApi) GetGitDiffShortStat(start, end string, opts ...Option) (*GitDiffShotStat, error) {
	res, err := g.GetDiffTwoCommit(start, end, opts...)
	if err != nil {
		return nil, err
	}

	sl := strings.Split(res, ",")
	if len(sl) == 3 {
		t := &GitDiffShotStat{}

		fc, err := strconv.Atoi(strings.Fields(sl[0])[0])
		if err == nil {
			t.FileChange = fc
		}
		in, err := strconv.Atoi(strings.Fields(sl[1])[0])
		if err == nil {
			t.InsertCount = in
		}
		dn, err := strconv.Atoi(strings.Fields(sl[2])[0])
		if err == nil {
			t.DeleteCount = dn
		}
		return t, nil
	} else if len(sl) == 2 {
		t := &GitDiffShotStat{}

		if strings.Contains(sl[0], "files changed") {
			dn, err := strconv.Atoi(strings.Fields(sl[0])[0])
			if err == nil {
				t.FileChange = dn
			}
		} else {
			dn, err := strconv.Atoi(strings.Fields(sl[0])[0])
			if err == nil {
				t.DeleteCount = dn
			}
		}

		if strings.Contains(sl[1], "insertions") {
			in, err := strconv.Atoi(strings.Fields(sl[1])[0])
			if err == nil {
				t.InsertCount = in
			}
		} else {
			dn, err := strconv.Atoi(strings.Fields(sl[1])[0])
			if err == nil {
				t.DeleteCount = dn
			}
		}
		return t, nil
	} else {
		return &GitDiffShotStat{}, err
	}

}

func (g *GitApi) GetGitDiffItemListStat(list []*GitDiffItem) GitDiffItemListStat {
	sort.Slice(list, func(i, j int) bool {
		if list[i].InsertNum != list[j].InsertNum {
			return list[i].InsertNum > list[j].InsertNum
		} else {
			if list[i].DeleteNum != list[j].DeleteNum {
				return list[i].DeleteNum > list[j].DeleteNum
			} else {
				return list[i].FileChangeNum > list[j].FileChangeNum
			}
		}
	})

	stat := GitDiffItemListStat{}
	for _, item := range list {
		stat.FileCount += item.FileChangeNum
		stat.InsertCount += item.InsertNum
		stat.DeleteCount += item.DeleteNum
	}
	return stat
}

func (g *GitApi) getDiffItemWithTargetExts(str string, exts map[string]string) *GitDiffItem {
	t := &GitDiffItem{
		Dir: g.Dir,
	}
	insertCount := 0
	deleteCount := 0
	fileCount := 0
	lines := stringx.SplitLines(str)
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
		m := lang_ext.CommonBackLanguageExt
		return g.getDiffItemWithTargetExts(str, m)
	} else if o.OnlyFrontLanguage {
		m := lang_ext.CommonFrontLanguageExt
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

// ---------------------------------------------------------------------------------------------------------------------

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
