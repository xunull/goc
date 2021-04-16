package git_cmd

import (
	"github.com/rs/zerolog/log"
	"github.com/xunull/goc/commandx"
	"github.com/xunull/goc/enhance/timex"
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

type DailyTwoLogItem struct {
	Dir   string
	Start time.Time
	End   time.Time
	Data  map[int]*TwoLogItem
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

func (g *GitApi) SortLogItemByDateAsc(items []LogItem) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Date.Before(items[j].Date)
	})
}

func (g *GitApi) SortRepoByCommitNum(logs []*RepoLogs) {
	sort.Slice(logs, func(i, j int) bool {
		return logs[i].Length > logs[j].Length
	})
}

func (g *GitApi) FormatLogRfc3339TimeItem(input string) *RepoLogs {
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

			t, err := time.Parse(time.RFC3339, date)
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

func (g *GitApi) FormatDefaultGitLogItem(input string) *RepoLogs {

	if g.option.PrettyFormat == PrettyRFC3339HashTime {
		return g.FormatLogRfc3339TimeItem(input)
	}

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

func (g *GitApi) GetLog(opts ...Option) (string, error) {
	o := g.getOption(opts...)
	g.option = o
	cmd := []string{"git", "--no-pager", "log"}

	if o.Since != "" {
		cmd = append(cmd, "--since="+o.Since)
	}

	if o.LogItemNum != 0 {
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

func (g *GitApi) GetLogSinceUntil(since, until time.Time, opts ...Option) (string, error) {
	o := g.getOption(opts...)
	g.option = o
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

func (g *GitApi) GetDailyTwoLog(logs *RepoLogs, start, end time.Time) *DailyTwoLogItem {

	days, err := timex.GetTwoDateDayIntList(start, end)
	if err != nil {
		return nil
	}

	m := make(map[int][]LogItem)
	for _, day := range days {
		m[day] = make([]LogItem, 0)
	}

	for _, targetLog := range logs.Logs {
		day, err := strconv.Atoi(targetLog.Date.Format("20060102"))
		if err != nil {
			return nil
		}
		m[day] = append(m[day], targetLog)
	}

	nm := make(map[int]*TwoLogItem)
	lastDay := 0

	for _, day := range days {
		if len(m[day]) > 0 {
			first := m[day][0]
			if lastDay > 0 {
				nm[lastDay].End = &first
			}
			nm[day] = &TwoLogItem{
				Start: &first,
			}
			lastDay = day
		}
	}

	return &DailyTwoLogItem{
		Dir:   logs.Dir,
		Start: start,
		End:   end,
		Data:  nm,
	}
}

// ---------------------------------------------------------------------------------------------------------------------

func (g *GitApi) FormatLogGraphOutput(out string) *RepoLogs {
	lines := strings.Split(out, "\n")
	if len(lines) > 0 {
		res := make([]LogItem, 0, 2)
		first := lines[0]
		fields := strings.Fields(first)

		lines = lines[1:]
		res = append(res, LogItem{
			Hash: fields[1],
			Dir:  g.Dir,
		})
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				t := reLogGraphBracket.Find([]byte(line))
				s := string(t)
				if s != "" {
					if !reOnlyTag.MatchString(s) {
						fields := strings.Fields(line)
						res = append(res, LogItem{
							Hash: fields[1],
							Dir:  g.Dir,
						})
					}
					return &RepoLogs{
						Logs:   res,
						Dir:    g.Dir,
						Length: len(res),
					}
				}
			}
		}
		return &RepoLogs{
			Logs:   nil,
			Dir:    g.Dir,
			Length: 0,
		}
	} else {
		return &RepoLogs{
			Logs:   nil,
			Dir:    g.Dir,
			Length: 0,
		}
	}

}

func (g *GitApi) GetCurrentLogGraph(opts ...Option) (string, error) {
	// todo use 1000
	opts = append(opts, WithLogItemLimit(1000))
	o := g.getOption(opts...)
	g.option = o
	cmd := []string{"git", "--no-pager", "log"}

	//cmd = append(cmd, "-n "+strconv.Itoa(o.LogItemNum))
	cmd = append(cmd, "--graph")
	cmd = append(cmd, "--oneline")
	cmd = append(cmd, "--decorate=short")
	cmd = append(cmd, "--no-color")

	cmdRes := commandx.RunCommand(cmd, commandx.WithDir(g.Dir))
	if cmdRes.Success {
		return strings.TrimSpace(cmdRes.Stdout.String()), nil
	} else {
		return cmdRes.Stderr.String(), cmdRes.Err
	}
}
