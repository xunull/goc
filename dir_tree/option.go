package dir_tree

type option struct {
	DefaultExclude  bool
	DotDirExclude   bool
	WithProgressBar bool
	SyncMode        bool
	SyncFileOpMode  bool
	TargetExt       string
	Depth           int
	OnlyDir         bool
	ExcludeSuffixes []string
	ExcludePrefixes []string
	ExcludeDir      []string
	excludeDirMap   map[string]struct{}
	ExcludeUnknown  bool
	WorkerCount     int
}

type Option func(o *option)

func getDefaultOption() *option {
	return &option{
		WorkerCount: 1024,
	}
}

func WithWorkerCount(count int) Option {
	return func(o *option) {
		o.WorkerCount = count
	}
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

func WithSyncFileOpMode() Option {
	return func(o *option) {
		o.SyncFileOpMode = true
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

func WithExcludeUnknown(exclude bool) Option {
	return func(o *option) {
		o.ExcludeUnknown = exclude
	}
}

func WithProgressBarOut() Option {
	return func(o *option) {
		o.WithProgressBar = true
	}
}
