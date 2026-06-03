package traverse_v2

import "runtime"

type Option func(*option)

type option struct {
	DirWorkers  int
	FileWorkers int
	QueueScale  int

	MaxDepth int

	TargetExt     string
	ExcludePrefix []string
	ExcludeSuffix []string
	ExcludeDir    []string
	excludeDirSet map[string]struct{}

	SkipDotEntries       bool
	SkipKnownIgnoreDirs  bool
	SkipKnownBinaryFiles bool

	OnlyDir bool

	OnDirComplete func(*Item)
	OnComplete    func()
}

func defaultOption() *option {
	return &option{
		DirWorkers:  64,
		FileWorkers: max(runtime.NumCPU(), 2) * 2,
		QueueScale:  16,
	}
}

func WithDirWorkers(n int) Option {
	return func(o *option) { o.DirWorkers = n }
}

func WithFileWorkers(n int) Option {
	return func(o *option) { o.FileWorkers = n }
}

// WithWorkerCount sets both DirWorkers and FileWorkers to n. v1
// compatibility shortcut; for separate tuning use WithDirWorkers and
// WithFileWorkers.
func WithWorkerCount(n int) Option {
	return func(o *option) {
		o.DirWorkers = n
		o.FileWorkers = n
	}
}

// WithSyncMode runs directory reads in a single worker. v1 compatibility
// alias for WithDirWorkers(1). Note: due to v2's overflow-on-full submit,
// this is "mostly single-threaded" — bursts may still spawn goroutines.
func WithSyncMode() Option {
	return func(o *option) { o.DirWorkers = 1 }
}

// WithSyncFileOpMode runs file callbacks in a single worker. v1
// compatibility alias for WithFileWorkers(1). Useful when the user
// callback is not goroutine-safe.
func WithSyncFileOpMode() Option {
	return func(o *option) { o.FileWorkers = 1 }
}

// WithQueueScale sets queue size = workers * scale for both pools. Default 16.
func WithQueueScale(n int) Option {
	return func(o *option) { o.QueueScale = n }
}

// WithMaxDepth limits traversal depth. 0 means unlimited.
func WithMaxDepth(depth int) Option {
	return func(o *option) { o.MaxDepth = depth }
}

// WithDepth is a v1 compatibility alias for WithMaxDepth.
func WithDepth(depth int) Option { return WithMaxDepth(depth) }

// WithTargetExt filters to only files matching this extension (callback is
// not invoked for other files). Pass with leading dot, e.g. ".go".
func WithTargetExt(ext string) Option {
	return func(o *option) { o.TargetExt = ext }
}

func WithExcludePrefix(prefixes ...string) Option {
	return func(o *option) { o.ExcludePrefix = append(o.ExcludePrefix, prefixes...) }
}

func WithExcludeSuffix(suffixes ...string) Option {
	return func(o *option) { o.ExcludeSuffix = append(o.ExcludeSuffix, suffixes...) }
}

// WithExcludeDir matches directory names (basename) OR relative paths.
// e.g. "node_modules" matches anywhere; "src/generated" matches that path.
func WithExcludeDir(dirs ...string) Option {
	return func(o *option) {
		o.ExcludeDir = append(o.ExcludeDir, dirs...)
		if o.excludeDirSet == nil {
			o.excludeDirSet = make(map[string]struct{}, len(o.ExcludeDir))
		}
		for _, d := range dirs {
			o.excludeDirSet[d] = struct{}{}
		}
	}
}

// WithSkipDotEntries skips files and directories whose name begins with ".".
func WithSkipDotEntries() Option {
	return func(o *option) { o.SkipDotEntries = true }
}

// WithSkipKnownIgnoreDirs skips directories matched by lang_ext.IsExcludeDir
// (node_modules, vendor, dist, etc.).
func WithSkipKnownIgnoreDirs() Option {
	return func(o *option) { o.SkipKnownIgnoreDirs = true }
}

// WithSkipKnownBinaryFiles skips files whose extension matches
// lang_ext.IsExcludeFileExt (.exe, .so, .pyc, etc.).
func WithSkipKnownBinaryFiles() Option {
	return func(o *option) { o.SkipKnownBinaryFiles = true }
}

// WithSensibleDefaults turns on all three common skip rules.
func WithSensibleDefaults() Option {
	return func(o *option) {
		o.SkipDotEntries = true
		o.SkipKnownIgnoreDirs = true
		o.SkipKnownBinaryFiles = true
	}
}

// WithDefaultExclude is a v1 compatibility alias for WithSensibleDefaults.
func WithDefaultExclude() Option { return WithSensibleDefaults() }

func WithOnlyDir() Option {
	return func(o *option) { o.OnlyDir = true }
}

// WithOnDirComplete registers a callback fired exactly once per directory,
// the moment that directory's entire subtree (all sub-dirs + all files)
// has completed. Useful for per-dir aggregation or incremental output.
// The root directory also fires this event last.
func WithOnDirComplete(fn func(*Item)) Option {
	return func(o *option) { o.OnDirComplete = fn }
}

// WithOnComplete registers a callback fired exactly once when the entire
// traversal has finished (after every item callback and every per-dir
// complete event). Fires before Done() unblocks and before Run() returns,
// so observers reading Done() see a fully-settled world.
func WithOnComplete(fn func()) Option {
	return func(o *option) { o.OnComplete = fn }
}
