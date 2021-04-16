package trend_prediction

const (
	DefaultWindowLength       = 3
	DefaultWindowEscapeLength = 60
)

type option struct {
	WindowLength     int
	WindowTimeEscape int
	OnlyUseLength    bool
	OnlyUseEscape    bool
	CheckIncreasing  bool
	CheckDescending  bool
	CheckAverage     bool
	AverageTarget    int
	Threshold        int
}

type Option func(o *option)

func getDefaultOption() *option {
	return &option{
		WindowLength:     DefaultWindowLength,
		WindowTimeEscape: DefaultWindowEscapeLength,
	}
}

func WithThreshold(threshold int) Option {
	return func(o *option) {
		o.Threshold = threshold
	}
}

func WithCheckAverage(target int) Option {
	return func(o *option) {
		o.CheckAverage = true
		o.AverageTarget = target
	}
}

func WithCheckIncreasing() Option {
	return func(o *option) {
		o.CheckIncreasing = true
	}
}

func WithCheckDescending() Option {
	return func(o *option) {
		o.CheckDescending = true
	}
}

func WithWindowLength(length int) Option {
	return func(o *option) {
		o.WindowLength = length
	}
}

func WithWindowTimeEscape(escape int) Option {
	return func(o *option) {
		o.WindowTimeEscape = escape
	}
}

func WithOnlyUseLength() Option {
	return func(o *option) {
		o.OnlyUseLength = true
	}
}

func WithOnlyUseEscape() Option {
	return func(o *option) {
		o.OnlyUseEscape = true
	}
}
