package commandx

import "time"

type option struct {
	Dir            string
	Timeout        time.Duration
	RedirectStderr bool
}

type Option func(o *option)

func WithDir(dir string) Option {
	return func(o *option) {
		o.Dir = dir
	}
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *option) {
		o.Timeout = timeout
	}
}

func WithRedirectStderr() Option {
	return func(o *option) {
		o.RedirectStderr = true
	}
}
