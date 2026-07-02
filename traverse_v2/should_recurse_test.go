package traverse_v2

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"sync"
	"testing"
)

// mustMkdirFile creates base/<rel> (with parent dirs) as an empty file.
func mustMkdirFile(t *testing.T, base, rel string) {
	t.Helper()
	full := filepath.Join(base, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, nil, 0644); err != nil {
		t.Fatal(err)
	}
}

func collectFilePaths(t *testing.T, root string, opts ...Option) []string {
	t.Helper()
	var mu sync.Mutex
	var files []string
	trv := New(root, func(item *Item) {
		if !item.IsDir {
			mu.Lock()
			files = append(files, item.Path)
			mu.Unlock()
		}
	}, opts...)
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	sort.Strings(files)
	return files
}

// A hook returning false for a directory must prune that entire subtree —
// none of its files or nested dirs are visited.
func TestWithShouldRecurse_PrunesSubtree(t *testing.T) {
	tmp := t.TempDir()
	mustMkdirFile(t, tmp, "keep/f.txt")
	mustMkdirFile(t, tmp, "skip/f.txt")
	mustMkdirFile(t, tmp, "skip/nested/g.txt")

	files := collectFilePaths(t, tmp, WithShouldRecurse(func(it *Item) bool {
		return !(it.IsDir && it.Name == "skip")
	}))

	want := []string{"keep/f.txt"}
	if !reflect.DeepEqual(files, want) {
		t.Errorf("expected %v, got %v", want, files)
	}
}

// A hook returning false for a file skips just that file, not its siblings.
func TestWithShouldRecurse_SkipsFile(t *testing.T) {
	tmp := t.TempDir()
	for _, n := range []string{"a.go", "b.txt", "c.go"} {
		if err := os.WriteFile(filepath.Join(tmp, n), nil, 0644); err != nil {
			t.Fatal(err)
		}
	}

	files := collectFilePaths(t, tmp, WithShouldRecurse(func(it *Item) bool {
		return it.IsDir || filepath.Ext(it.Name) == ".go"
	}))

	want := []string{"a.go", "c.go"}
	if !reflect.DeepEqual(files, want) {
		t.Errorf("expected %v, got %v", want, files)
	}
}

// The hook Item must carry the rel Path so callers (e.g. a gitignore matcher)
// can decide based on the path, not just the basename.
func TestWithShouldRecurse_ItemCarriesRelPath(t *testing.T) {
	tmp := t.TempDir()
	mustMkdirFile(t, tmp, "a/b/c.txt")

	var mu sync.Mutex
	seen := map[string]bool{}
	trv := New(tmp, func(*Item) {}, WithShouldRecurse(func(it *Item) bool {
		mu.Lock()
		seen[it.Path] = true
		mu.Unlock()
		return true
	}))
	if err := trv.Run(); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{"a", "a/b", "a/b/c.txt"} {
		if !seen[p] {
			t.Errorf("hook was not called with rel Path %q; saw %v", p, seen)
		}
	}
}

// CRITICAL regression: a permissive hook (always true) and no hook at all must
// produce EXACTLY the same traversal as the pre-hook default. This guards
// du/cloc/clocx/find/grep against behavior drift from adding the hook.
func TestWithShouldRecurse_PermissiveMatchesDefault(t *testing.T) {
	tmp := t.TempDir()
	var td, tf int32
	if err := buildTree(tmp, 4, 3, 3, &td, &tf); err != nil {
		t.Fatal(err)
	}

	base := collectFilePaths(t, tmp)
	withTrueHook := collectFilePaths(t, tmp, WithShouldRecurse(func(*Item) bool { return true }))

	if !reflect.DeepEqual(base, withTrueHook) {
		t.Errorf("permissive hook changed traversal: base=%d items, hook=%d items",
			len(base), len(withTrueHook))
	}
}

// The hook is ANDed with static excludes: if either says skip, the entry is
// skipped. A permissive hook must not override WithSkipKnownIgnoreDirs.
func TestWithShouldRecurse_ANDsWithStaticExclude(t *testing.T) {
	tmp := t.TempDir()
	mustMkdirFile(t, tmp, "node_modules/f.txt")
	mustMkdirFile(t, tmp, "keep/f.txt")

	files := collectFilePaths(t, tmp,
		WithSkipKnownIgnoreDirs(),
		WithShouldRecurse(func(*Item) bool { return true }),
	)

	want := []string{"keep/f.txt"}
	if !reflect.DeepEqual(files, want) {
		t.Errorf("expected %v (node_modules pruned by static exclude), got %v", want, files)
	}
}
