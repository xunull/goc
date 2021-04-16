package git_cmd

import "fmt"

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

func WithSince(day uint) Option {
	return func(o *option) {
		o.Since = fmt.Sprintf("'%d day ago'", day)
	}
}

func (g *GitApi) getOption(opts ...Option) *option {
	d := &option{
		LogItemNum:   20,
		PrettyFormat: "format:%h %cs",
	}
	for _, o := range opts {
		o(d)
	}
	return d
}
