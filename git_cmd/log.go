package git_cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/xunull/goc/commandx"
	"sort"
	"strconv"
	"strings"
	"time"
)

type LogItem struct {
	Hash string
	Date time.Time
	Dir  string
}

type TwoLogItem struct {
	Start *LogItem
	End   *LogItem
}

type RepoLogs struct {
	Dir           string
	Logs          []LogItem
	Length        int
	FromFirstDiff *GitDiffItem
}

func (g *GitApi) DiffFromFirst(rl *RepoLogs) (*GitDiffItem, error) {
	g.SortLogItemByDateDesc(rl.Logs)
	first := rl.Logs[rl.Length-1]
	res, err := g.GetDiffCommitAndWorkTree(first.Hash, WithShortStat())
	if err != nil {
		return nil, err
	}
	dfItem := g.GetGitDiffItem(res)
	if dfItem != nil {
		dfItem.Dir = g.Dir
		return dfItem, nil
	}
	return nil, nil
}

func (g *GitApi) SortRepoLogByFirstDiff(items []*RepoLogs) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].FromFirstDiff.InsertNum-
			items[i].FromFirstDiff.DeleteNum >
			items[j].FromFirstDiff.InsertNum-
				items[j].FromFirstDiff.DeleteNum
	})
}

func (g *GitApi) SortLogItemByDateDesc(items []LogItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.After(items[j].Date)
	})
}

func (g *GitApi) SortRepoByCommitNum(logs []*RepoLogs) {
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Length > logs[j].Length
	})
}

func (g *GitApi) FormatDefaultGitLogItem(input string) *RepoLogs {
	if input == "" {
		return nil
	}
	lines := strings.Split(input, "\n")
	res := make([]LogItem, 0, len(lines))
	for _, line := range lines {
		a := strings.Fields(line)
		if len(a) == 2 {
			hash := a[0]
			date := a[1]

			t, err := time.Parse("2006-01-02", date)
			if err != nil {
				log.Info().Err(err).Msg("parse git log item date failed")
			}
			res = append(res, LogItem{
				Hash: hash,
				Dir:  g.Dir,
				Date: t,
			})
		}
	}
	return &RepoLogs{
		Logs:   res,
		Dir:    g.Dir,
		Length: len(res),
	}
}

func (g *GitApi) GetParent(hash string) (string, error) {
	cmd := []string{"git", "log", "-1", "--pretty=%p", hash}
	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) GetLog(ops ...Option) (string, error) {
	o := g.getOption(ops...)
	cmd := []string{"git", "log"}

	if o.Since != "" {
		cmd = append(cmd, "--since="+o.Since)
	} else {
		cmd = append(cmd, "-"+strconv.Itoa(o.LogItemNum))
	}

	cmd = append(cmd, "--pretty="+o.PrettyFormat)

	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}

func (g *GitApi) GetLogSinceUntil(since, until time.Time) (string, error) {
	o := g.getOption()
	cmd := []string{"git", "log"}

	cmd = append(cmd, "--since="+since.Format("2006-01-02"))

	cmd = append(cmd, "--until="+until.Format("2006-01-02"))

	cmd = append(cmd, "--pretty="+o.PrettyFormat)

	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}
