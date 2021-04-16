package file_path

type option struct {
	Suffix string
}

type Option func(o *option)

func getOption(opts ...Option) *option {
	op := &option{}
	for _, f := range opts {
		f(op)
	}
	return op
}

func WithSuffix(suffix string) Option {
	return func(o *option) {
		o.Suffix = suffix
	}
}
