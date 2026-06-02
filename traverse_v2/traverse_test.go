package traverse_v2

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestBasic(t *testing.T) {
	tmp := t.TempDir()
	for i := 0; i < 5; i++ {
		if err := os.WriteFile(filepath.Join(tmp, fmt.Sprintf("f%d.txt", i)), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	var seen atomic.Int32
	trv := New(tmp, func(item *Item) {
		if !item.IsDir {
			seen.Add(1)
		}
	})
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	if got := seen.Load(); got != 5 {
		t.Errorf("expected 5 files, got %d", got)
	}
}

func TestEmptyDir(t *testing.T) {
	tmp := t.TempDir()
	trv := New(tmp, func(item *Item) {})
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
}

func TestDeepNested(t *testing.T) {
	tmp := t.TempDir()
	var totalFiles, totalDirs int32
	if err := buildTree(tmp, 5, 3, 3, &totalDirs, &totalFiles); err != nil {
		t.Fatal(err)
	}

	var sawFiles, sawDirs atomic.Int32
	trv := New(tmp, func(item *Item) {
		if item.IsDir {
			sawDirs.Add(1)
		} else {
			sawFiles.Add(1)
		}
	})
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	if sawFiles.Load() != totalFiles {
		t.Errorf("expected %d files, got %d", totalFiles, sawFiles.Load())
	}
	if sawDirs.Load() != totalDirs {
		t.Errorf("expected %d dirs, got %d", totalDirs, sawDirs.Load())
	}
}

// TestRaceCompletion is the regression test for the WorkSheet race described
// in the v1 review: Run() must not return until every callback has finished.
// Wide-shallow tree with small files maximizes the chance of fast callbacks
// interleaving sibling Adds, the worst case for the v1 design.
func TestRaceCompletion(t *testing.T) {
	tmp := t.TempDir()
	const N = 500
	for i := 0; i < N; i++ {
		d := filepath.Join(tmp, fmt.Sprintf("d%d", i))
		if err := os.Mkdir(d, 0755); err != nil {
			t.Fatal(err)
		}
		for j := 0; j < 3; j++ {
			if err := os.WriteFile(filepath.Join(d, fmt.Sprintf("f%d", j)), nil, 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	want := int32(N + N*3) // N dirs + N*3 files
	for iter := 0; iter < 30; iter++ {
		var processed atomic.Int32
		trv := New(tmp, func(item *Item) {
			processed.Add(1)
		})
		if err := trv.Run(); err != nil {
			t.Fatalf("iter %d: %v", iter, err)
		}
		if got := processed.Load(); got != want {
			t.Fatalf("iter %d: Run() returned but %d/%d items processed", iter, got, want)
		}
	}
}

func TestDirCompleteEvent(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	if err := os.Mkdir(sub, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sub, "f.txt"), nil, 0644); err != nil {
		t.Fatal(err)
	}

	var dirComplete atomic.Int32
	var rootDoneLast atomic.Bool
	trv := New(tmp,
		func(item *Item) {},
		WithOnDirComplete(func(item *Item) {
			n := dirComplete.Add(1)
			// Root has empty Path; it must fire last.
			if item.Path == "" && n == 2 {
				rootDoneLast.Store(true)
			}
		}),
	)
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	if got := dirComplete.Load(); got != 2 {
		t.Errorf("expected 2 dir-complete events (sub + root), got %d", got)
	}
	if !rootDoneLast.Load() {
		t.Errorf("expected root dir-complete event to fire last")
	}
}

func TestOnCompleteFiresOnce(t *testing.T) {
	tmp := t.TempDir()
	for i := 0; i < 10; i++ {
		if err := os.WriteFile(filepath.Join(tmp, fmt.Sprintf("f%d", i)), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	var fired atomic.Int32
	trv := New(tmp, func(item *Item) {}, WithOnComplete(func() {
		fired.Add(1)
	}))
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	if got := fired.Load(); got != 1 {
		t.Errorf("expected OnComplete to fire exactly once, got %d", got)
	}
}

func TestDoneChannel(t *testing.T) {
	tmp := t.TempDir()
	for i := 0; i < 50; i++ {
		if err := os.WriteFile(filepath.Join(tmp, fmt.Sprintf("f%d", i)), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	trv := New(tmp, func(item *Item) {})

	// Before Run(): Done() must not be closed yet.
	select {
	case <-trv.Done():
		t.Fatalf("Done() closed before Run()")
	default:
	}

	go func() { _ = trv.Run() }()

	// After Run() starts, Done() should close within a reasonable time.
	select {
	case <-trv.Done():
		// ok
	case <-time.After(5 * time.Second):
		t.Fatalf("Done() did not close within 5s")
	}
}

// TestSignalOrdering verifies the documented order:
//   OnDirComplete(root) → OnComplete → Done() closed → Run() returns
func TestSignalOrdering(t *testing.T) {
	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "f"), nil, 0644); err != nil {
		t.Fatal(err)
	}

	var (
		mu              sync.Mutex
		order           []string
		seenRootDirDone bool
		seenOnComplete  bool
	)

	// Declare trv first so the WithOnComplete closure can reference it.
	var trv *Traverse
	trv = New(tmp, func(item *Item) {},
		WithOnDirComplete(func(item *Item) {
			if item.Path == "" {
				mu.Lock()
				order = append(order, "rootDirComplete")
				seenRootDirDone = true
				mu.Unlock()
			}
		}),
		WithOnComplete(func() {
			mu.Lock()
			order = append(order, "onComplete")
			seenOnComplete = true
			// Done() must NOT be closed yet at the moment OnComplete fires.
			select {
			case <-trv.Done():
				t.Errorf("Done() already closed when OnComplete fired")
			default:
			}
			mu.Unlock()
		}),
	)

	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}

	// After Run() returns: Done() must be closed.
	select {
	case <-trv.Done():
	default:
		t.Errorf("Done() not closed after Run() returned")
	}

	if !seenRootDirDone || !seenOnComplete {
		t.Fatalf("missing signals: rootDir=%v onComplete=%v", seenRootDirDone, seenOnComplete)
	}
	if len(order) != 2 || order[0] != "rootDirComplete" || order[1] != "onComplete" {
		t.Errorf("wrong order: %v", order)
	}
}

func TestTargetExt(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"a.go", "b.go", "c.txt", "d.md"} {
		if err := os.WriteFile(filepath.Join(tmp, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
	var n atomic.Int32
	trv := New(tmp, func(item *Item) {
		if !item.IsDir {
			n.Add(1)
		}
	}, WithTargetExt(".go"))
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	if got := n.Load(); got != 2 {
		t.Errorf("expected 2 .go files, got %d", got)
	}
}

func TestExcludeDir(t *testing.T) {
	tmp := t.TempDir()
	for _, d := range []string{"keep", "node_modules", "skipme"} {
		dp := filepath.Join(tmp, d)
		if err := os.Mkdir(dp, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dp, "f.txt"), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	var n atomic.Int32
	trv := New(tmp, func(item *Item) {
		if !item.IsDir {
			n.Add(1)
		}
	},
		WithSkipKnownIgnoreDirs(),
		WithExcludeDir("skipme"),
	)
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	if got := n.Load(); got != 1 { // only "keep/f.txt"
		t.Errorf("expected 1 file, got %d", got)
	}
}

func TestGetAllPaths(t *testing.T) {
	tmp := t.TempDir()
	if err := os.Mkdir(filepath.Join(tmp, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"a.txt", "b.txt", "sub/c.txt"} {
		full := filepath.Join(tmp, filepath.FromSlash(p))
		if err := os.WriteFile(full, nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	paths, err := GetAllPaths(tmp)
	if err != nil {
		t.Fatal(err)
	}
	sort.Strings(paths)
	want := []string{"a.txt", "b.txt", "sub/c.txt"}
	if len(paths) != len(want) {
		t.Fatalf("got %d paths, want %d: %v", len(paths), len(want), paths)
	}
	for i, p := range want {
		if paths[i] != p {
			t.Errorf("paths[%d] = %q, want %q", i, paths[i], p)
		}
	}
}

func TestGetFileCount(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"a.go", "b.go", "c.py", "d.txt", "no_ext_file"} {
		if err := os.WriteFile(filepath.Join(tmp, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	stats, err := GetFileCount(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if stats.Total != 5 {
		t.Errorf("Total = %d, want 5", stats.Total)
	}
	if stats.ByLanguage["Golang"] != 2 {
		t.Errorf("ByLanguage[Golang] = %d, want 2", stats.ByLanguage["Golang"])
	}
	if stats.ByLanguage["Python"] != 1 {
		t.Errorf("ByLanguage[Python] = %d, want 1", stats.ByLanguage["Python"])
	}
	if stats.ByLanguage["Text"] != 1 {
		t.Errorf("ByLanguage[Text] = %d, want 1", stats.ByLanguage["Text"])
	}
	// no_ext_file should be in Total but not in any language bucket.
}

func TestGetFileList(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"a.go", "b.txt", "c.go"} {
		if err := os.WriteFile(filepath.Join(tmp, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	res, err := GetFileList(tmp, WithTargetExt(".go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.List) != 2 {
		t.Fatalf("List len = %d, want 2: %v", len(res.List), res.List)
	}
	if _, ok := res.Set["a.go"]; !ok {
		t.Errorf("Set missing a.go")
	}
	if _, ok := res.Set["c.go"]; !ok {
		t.Errorf("Set missing c.go")
	}
	if _, ok := res.Set["b.txt"]; ok {
		t.Errorf("Set should not contain b.txt")
	}
}

// TestGetFileList_NoTargetExt confirms v2 fixes the v1 bug where an empty
// list was returned when WithTargetExt was not set.
func TestGetFileList_NoTargetExt(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"a", "b", "c"} {
		if err := os.WriteFile(filepath.Join(tmp, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}
	res, err := GetFileList(tmp)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.List) != 3 {
		t.Errorf("List len = %d, want 3 (v1 bug: returned 0)", len(res.List))
	}
}

func TestV1AliasOptions(t *testing.T) {
	// Verify the v1 compatibility aliases actually take effect by
	// constructing options and checking they set the underlying fields.
	o := defaultOption()
	WithWorkerCount(7)(o)
	if o.DirWorkers != 7 || o.FileWorkers != 7 {
		t.Errorf("WithWorkerCount(7): DirWorkers=%d FileWorkers=%d", o.DirWorkers, o.FileWorkers)
	}

	o = defaultOption()
	WithSyncMode()(o)
	if o.DirWorkers != 1 {
		t.Errorf("WithSyncMode: DirWorkers=%d, want 1", o.DirWorkers)
	}

	o = defaultOption()
	WithSyncFileOpMode()(o)
	if o.FileWorkers != 1 {
		t.Errorf("WithSyncFileOpMode: FileWorkers=%d, want 1", o.FileWorkers)
	}

	o = defaultOption()
	WithDepth(5)(o)
	if o.MaxDepth != 5 {
		t.Errorf("WithDepth(5): MaxDepth=%d, want 5", o.MaxDepth)
	}

	o = defaultOption()
	WithDefaultExclude()(o)
	if !o.SkipDotEntries || !o.SkipKnownIgnoreDirs || !o.SkipKnownBinaryFiles {
		t.Errorf("WithDefaultExclude: skip flags not all true: %+v", o)
	}
}

func TestGetAllPathsTargetExt(t *testing.T) {
	tmp := t.TempDir()
	for _, name := range []string{"a.go", "b.go", "c.txt", "d.md"} {
		if err := os.WriteFile(filepath.Join(tmp, name), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	paths, err := GetAllPaths(tmp, WithTargetExt(".go"))
	if err != nil {
		t.Fatal(err)
	}
	if len(paths) != 2 {
		t.Fatalf("got %d, want 2: %v", len(paths), paths)
	}
}

func buildTree(root string, depth, branchDir, branchFile int, dirCount, fileCount *int32) error {
	if depth == 0 {
		return nil
	}
	for i := 0; i < branchFile; i++ {
		if err := os.WriteFile(filepath.Join(root, fmt.Sprintf("f%d.txt", i)), nil, 0644); err != nil {
			return err
		}
		atomic.AddInt32(fileCount, 1)
	}
	for i := 0; i < branchDir; i++ {
		d := filepath.Join(root, fmt.Sprintf("d%d", i))
		if err := os.Mkdir(d, 0755); err != nil {
			return err
		}
		atomic.AddInt32(dirCount, 1)
		if err := buildTree(d, depth-1, branchDir, branchFile, dirCount, fileCount); err != nil {
			return err
		}
	}
	return nil
}
