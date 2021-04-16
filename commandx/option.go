package commandx

type option struct {
	Dir string
}

type Option func(o *option)

func WithDir(dir string) Option {
	return func(o *option) {
		o.Dir = dir
	}
}
