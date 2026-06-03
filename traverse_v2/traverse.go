// Package traverse_v2 is a rewrite of traverse aimed at maximum speed and
// correct hierarchical completion tracking.
//
// Design summary:
//   - Each directory owns a dirNode "signing sheet" with atomic
//     expected/done counters and a one-shot finished CAS. Completion
//     bubbles up the parent chain; root completion closes Traverse.done.
//   - Two independent worker pools: one for directory reads (IO-bound),
//     one for file callbacks (user code). Submit is non-blocking with
//     overflow-goroutine fallback to prevent recursive-submission deadlocks.
//   - Directory entries are read via os.ReadDir + entry.Type(), which
//     avoids the per-entry lstat syscall that os.DirEntry.Info() incurs.
package traverse_v2

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xunull/goc/lang_ext"
)

type Traverse struct {
	Path string

	opt    *option
	onItem func(*Item)

	dirPool  *pool
	filePool *pool

	done   chan struct{}
	ctx    context.Context
	cancel context.CancelFunc

	errMu sync.Mutex
	errs  []error
}

func New(path string, onItem func(*Item), opts ...Option) *Traverse {
	o := defaultOption()
	for _, fn := range opts {
		fn(o)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Traverse{
		Path:   path,
		opt:    o,
		onItem: onItem,
		ctx:    ctx,
		cancel: cancel,
		done:   make(chan struct{}),
	}
}

// Run blocks until the entire subtree under Path has been traversed and all
// file callbacks have returned. Errors during ReadDir are collected and
// reported via the returned error (first one wrapped) and Errors().
func (t *Traverse) Run() error {
	t.dirPool = newPool(t.opt.DirWorkers, t.opt.DirWorkers*t.opt.QueueScale)
	t.filePool = newPool(t.opt.FileWorkers, t.opt.FileWorkers*t.opt.QueueScale)
	defer t.dirPool.close()
	defer t.filePool.close()

	root := &dirNode{
		t:    t,
		path: t.Path,
		rel:  "",
		name: filepath.Base(t.Path),
	}

	t.dirPool.submit(func() { t.traverseDir(root) })
	<-t.done

	return t.firstError()
}

func (t *Traverse) Cancel() { t.cancel() }

// Done returns a channel that is closed when the entire traversal has
// finished — every item callback returned, every per-dir complete event
// fired, and the OnComplete callback (if any) finished. Safe to call
// before Run(); the channel is created at construction time.
//
// Use it from a goroutine that does NOT call Run() itself, e.g.:
//
//	go trv.Run()
//	<-trv.Done()
func (t *Traverse) Done() <-chan struct{} { return t.done }

func (t *Traverse) Errors() []error {
	t.errMu.Lock()
	defer t.errMu.Unlock()
	return append([]error(nil), t.errs...)
}

func (t *Traverse) HasErrors() bool {
	t.errMu.Lock()
	defer t.errMu.Unlock()
	return len(t.errs) > 0
}

func (t *Traverse) recordErr(err error) {
	t.errMu.Lock()
	t.errs = append(t.errs, err)
	t.errMu.Unlock()
}

func (t *Traverse) firstError() error {
	t.errMu.Lock()
	defer t.errMu.Unlock()
	if len(t.errs) == 0 {
		return nil
	}
	return fmt.Errorf("traverse completed with %d errors: first error: %w", len(t.errs), t.errs[0])
}

func (t *Traverse) traverseDir(n *dirNode) {
	// Cancellation: still must close discovery so parent can finalize.
	select {
	case <-t.ctx.Done():
		n.closeDiscovery()
		return
	default:
	}

	if t.opt.MaxDepth > 0 && n.depth >= t.opt.MaxDepth {
		n.closeDiscovery()
		return
	}

	// Fire onItem for this directory (skip root to match v1 semantics where
	// the root path isn't reported as an item).
	if n.parent != nil && t.onItem != nil {
		t.onItem(&Item{
			Path:     n.rel,
			FullPath: n.path,
			Name:     n.name,
			IsDir:    true,
			Depth:    n.depth,
			Mode:     fs.ModeDir,
		})
	}

	entries, err := os.ReadDir(n.path)
	if err != nil {
		t.recordErr(err)
		n.closeDiscovery()
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		mode := entry.Type()

		if mode&fs.ModeSymlink != 0 {
			continue
		}

		if entry.IsDir() {
			if !t.shouldRecurseDir(name, n.rel) {
				continue
			}
			n.addChild()
			sub := &dirNode{
				t:      t,
				parent: n,
				path:   filepath.Join(n.path, name),
				rel:    joinRel(n.rel, name),
				name:   name,
				depth:  n.depth + 1,
			}
			t.dirPool.submit(func() { t.traverseDir(sub) })
			continue
		}

		ext := filepath.Ext(name)
		if !t.shouldProcessFile(name, ext) {
			continue
		}

		n.addChild()
		item := &Item{
			Path:     joinRel(n.rel, name),
			FullPath: filepath.Join(n.path, name),
			Name:     name,
			Ext:      ext,
			Mode:     mode,
			IsDir:    false,
			Depth:    n.depth + 1,
		}
		parent := n
		t.filePool.submit(func() {
			if t.onItem != nil {
				t.onItem(item)
			}
			parent.childDone()
		})
	}

	n.closeDiscovery()
}

func (t *Traverse) shouldRecurseDir(name, parentRel string) bool {
	if t.opt.SkipDotEntries && strings.HasPrefix(name, ".") {
		return false
	}
	if t.opt.SkipKnownIgnoreDirs {
		if lang_ext.IsExcludeDir(name) {
			return false
		}
	}
	if t.opt.excludeDirSet != nil {
		if _, ok := t.opt.excludeDirSet[name]; ok {
			return false
		}
		rel := joinRel(parentRel, name)
		if _, ok := t.opt.excludeDirSet[rel]; ok {
			return false
		}
	}
	return true
}

func (t *Traverse) shouldProcessFile(name, ext string) bool {
	if t.opt.OnlyDir {
		return false
	}
	if t.opt.SkipDotEntries && strings.HasPrefix(name, ".") {
		return false
	}
	if t.opt.SkipKnownBinaryFiles {
		if lang_ext.IsExcludeFileExt(ext) {
			return false
		}
	}
	if t.opt.TargetExt != "" && ext != t.opt.TargetExt {
		return false
	}
	for _, sfx := range t.opt.ExcludeSuffix {
		if strings.HasSuffix(name, sfx) {
			return false
		}
	}
	for _, pfx := range t.opt.ExcludePrefix {
		if strings.HasPrefix(name, pfx) {
			return false
		}
	}
	return true
}

// joinRel concatenates two path components with "/" without invoking
// filepath.Join (which Clean()s and allocates more). Hot path.
func joinRel(parent, name string) string {
	if parent == "" {
		return name
	}
	return parent + "/" + name
}
