package traverse_v2

import (
	"sync"

	"github.com/xunull/goc/lang_ext"
)

// GetAllPaths walks dir and returns the relative paths of every file
// (directories are not included). Paths use "/" as separator and are
// rooted at dir (e.g. "subdir/file.go", not "/abs/path/subdir/file.go").
//
// All Option values are honored. Common usage:
//
//	paths, err := traverse_v2.GetAllPaths(root,
//	    traverse_v2.WithSensibleDefaults(),
//	    traverse_v2.WithTargetExt(".go"),
//	)
//
// This is the v2 equivalent of v1 DirTraverse.GetAllPath.
func GetAllPaths(dir string, opts ...Option) ([]string, error) {
	var (
		mu    sync.Mutex
		paths []string
	)
	trv := New(dir, func(item *Item) {
		if item.IsDir {
			return
		}
		mu.Lock()
		paths = append(paths, item.Path)
		mu.Unlock()
	}, opts...)
	err := trv.Run()
	return paths, err
}

// FileCountResult is what GetFileCount returns. v2 equivalent of v1
// traverse.FileCountRes, minus TargetCount (redundant — when TargetExt is
// set, Total already equals matching files) and minus the embedded option.
type FileCountResult struct {
	// Total file count after filtering (directories not counted, unlike
	// v1 which over-counted by including dirs in Count).
	Total int

	// File count grouped by language name, keyed via lang_ext.CommonLanguageExt.
	// Files with unrecognized extensions are counted in Total but not here.
	ByLanguage map[string]int
}

// GetFileCount walks dir and returns a count of files plus per-language
// breakdown. v2 equivalent of v1 traverse.GetFileCount.
func GetFileCount(dir string, opts ...Option) (*FileCountResult, error) {
	var mu sync.Mutex
	res := &FileCountResult{ByLanguage: make(map[string]int)}

	trv := New(dir, func(item *Item) {
		if item.IsDir {
			return
		}
		mu.Lock()
		res.Total++
		if lang, ok := lang_ext.CommonLanguageExt[item.Ext]; ok {
			res.ByLanguage[lang]++
		}
		mu.Unlock()
	}, opts...)

	err := trv.Run()
	return res, err
}

// FileListResult bundles a file path slice with an O(1) lookup set. v2
// equivalent of v1 traverse.FileListRes.
type FileListResult struct {
	List []string
	Set  map[string]struct{}
}

// GetFileList walks dir and returns matching file paths (relative to dir,
// "/" separated) plus an O(1) lookup set. Combine with WithTargetExt to
// filter by extension.
//
// Difference from v1: when TargetExt is NOT set, v2 includes ALL files
// (v1 returned an empty list — a bug in v1's callback that this fixes).
func GetFileList(dir string, opts ...Option) (*FileListResult, error) {
	var mu sync.Mutex
	res := &FileListResult{
		List: nil,
		Set:  make(map[string]struct{}),
	}

	trv := New(dir, func(item *Item) {
		if item.IsDir {
			return
		}
		mu.Lock()
		res.List = append(res.List, item.Path)
		res.Set[item.Path] = struct{}{}
		mu.Unlock()
	}, opts...)

	err := trv.Run()
	return res, err
}
