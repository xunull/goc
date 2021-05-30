package boltkv_cache

import "time"

type option struct {
	CreateAt       time.Time
	VisitAt        time.Time
	ExpireAt       time.Time
	ExpireDuration time.Duration
}

type Option func(o *option)

func getOption(opts ...Option) *option {
	o := &option{}
	for _, f := range opts {
		f(o)
	}
	return o
}

func WithDuration(d time.Duration) Option {
	return func(o *option) {
		o.ExpireDuration = d
	}
}

func WithExpireAt(t time.Time) Option {
	return func(o *option) {
		o.ExpireAt = t
	}
}
