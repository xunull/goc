package git_cmd

import "fmt"

const (
	PrettyRFC3339HashTime = "format:%h %cI"
)

type option struct {
	ShortStat         bool
	PrettyFormat      string
	LogItemNum        int
	Since             string
	AddAll            bool
	SinceDay          string
	UntilDay          string
	TargetExt         string
	ExcludeExt        []string
	OnlyFrontLanguage bool
	OnlyBackLanguage  bool
	LogNoPage         bool
}

type Option func(o *option)

func WithOnlyFrontLanguage() Option {
	return func(o *option) {
		o.OnlyFrontLanguage = true
	}
}

func WithOnlyBackLanguage() Option {
	return func(o *option) {
		o.OnlyBackLanguage = true
	}
}

func WithExcludeExt(exclude []string) Option {
	return func(o *option) {
		o.ExcludeExt = exclude
	}
}

func WithLogNoPage() Option {
	return func(o *option) {
		o.LogNoPage = true
	}
}

func WithTargetExt(ext string) Option {
	return func(o *option) {
		o.TargetExt = ext
	}
}

func WithSinceDay(day string) Option {
	return func(o *option) {
		o.SinceDay = day
	}
}

func WithUntilDay(day string) Option {
	return func(o *option) {
		o.UntilDay = day
	}
}

func WithAddAll() Option {
	return func(o *option) {
		o.AddAll = true
	}
}

func WithShortStat() Option {
	return func(o *option) {
		o.ShortStat = true
	}
}

func WithNumStat() Option {
	return func(o *option) {
		o.ShortStat = false
	}
}

func WithLogPrettyRFC3339() Option {
	return func(o *option) {
		o.PrettyFormat = PrettyRFC3339HashTime
	}
}

func WithDefaultPrettyFormat() Option {
	return func(o *option) {
		o.PrettyFormat = "format:%h %cs"
	}
}

func WithDefaultLogItemNum() Option {
	return func(o *option) {
		o.LogItemNum = 20
	}
}

func WithLogItemLimit(limit int) Option {
	return func(o *option) {
		o.LogItemNum = limit
	}
}

func WithSince(day uint) Option {
	return func(o *option) {
		o.Since = fmt.Sprintf("'%d day ago'", day)
	}
}

func (g *GitApi) getOption(opts ...Option) *option {
	d := &option{
		LogItemNum:   0,
		PrettyFormat: "format:%h %cs",
	}
	for _, o := range opts {
		o(d)
	}
	return d
}
