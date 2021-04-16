package tree

type Option struct {
	IgnoreDefault     bool
	TraverseWithCache bool
}

type TreeOption func(o *Option)

func GetDefaultOption() *Option {
	return &Option{
		IgnoreDefault: true,
	}
}

func getTreeOption(opts ...TreeOption) *Option {
	t := &Option{}
	for _, o := range opts {
		o(t)
	}
	return t
}

func WithTraverseCache() TreeOption {
	return func(o *Option) {
		o.TraverseWithCache = true
	}
}
