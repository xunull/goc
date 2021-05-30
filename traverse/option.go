package traverse

type option struct {
	DefaultExclude  bool
	DotDirExclude   bool
	WithProgressBar bool
	SyncMode        bool
	TargetExt       string
	Depth           int
	OnlyDir         bool
	ExcludeSuffixes []string
	ExcludePrefixes []string
	ExcludeDir      []string
	excludeDirMap   map[string]struct{}
}

type Option func(o *option)

func getDefaultOption() *option {
	return &option{}
}

func WithExcludeDir(exclude []string) Option {
	return func(o *option) {
		o.ExcludeDir = exclude
		o.excludeDirMap = make(map[string]struct{})
		for _, dir := range exclude {
			o.excludeDirMap[dir] = struct{}{}
		}
	}
}

func WithExcludeSuffix(list ...string) Option {
	return func(o *option) {
		if o.ExcludeSuffixes == nil {
			o.ExcludeSuffixes = make([]string, 0)
		}
		o.ExcludeSuffixes = append(o.ExcludeSuffixes, list...)
	}
}

func WithExcludePrefix(list ...string) Option {
	return func(o *option) {
		if o.ExcludePrefixes == nil {
			o.ExcludePrefixes = make([]string, 0)
		}
		o.ExcludePrefixes = append(o.ExcludePrefixes, list...)
	}
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
