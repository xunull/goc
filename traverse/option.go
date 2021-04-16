package traverse

type option struct {
	DefaultExclude  bool
	DotDirExclude   bool
	WithProgressBar bool
	SyncMode        bool
	TargetExt       string
	Depth           int
	OnlyDir         bool
}

type Option func(o *option)

func getDefaultOption() *option {
	return &option{}
}

func WithOnlyDir() Option {
	return func(o *option) {
		o.OnlyDir = true
	}
}

func WithDepth(depth int) Option {
	return func(o *option) {
		o.Depth = depth
	}
}

func WithSyncMode() Option {
	return func(o *option) {
		o.SyncMode = true
	}
}

func WithTargetExt(ext string) Option {
	return func(o *option) {
		o.TargetExt = ext
	}
}

func WithDefaultExclude() Option {
	return func(o *option) {
		o.DefaultExclude = true
		o.DotDirExclude = true
	}
}

func WithProgressBarOut() Option {
	return func(o *option) {
		o.WithProgressBar = true
	}
}
