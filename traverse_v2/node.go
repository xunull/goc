package traverse_v2

import "sync/atomic"

// dirNode is the per-directory signing sheet (签单). Each directory owns
// one. The four atomics cooperate to drive completion bubbling without
// any mutex or channel.
//
// Invariants:
//   - expected only grows, and only by the parent goroutine running this
//     node's traverseDir for-loop (single writer).
//   - done only grows, written from many goroutines via atomic.Add.
//   - discComplete is set exactly once, after the for-loop finishes.
//   - finished is a one-shot CAS guard ensuring completion fires at most once.
//
// A node completes iff (discComplete == true) AND (done >= expected). The
// discComplete gate is what prevents the "Count==OverCount transient race"
// that the v1 WorkSheet has: even if done == expected mid-traversal, we
// refuse to finalize until expected is closed.
type dirNode struct {
	t      *Traverse
	parent *dirNode

	path  string
	rel   string
	name  string
	depth int

	expected     atomic.Int32
	done         atomic.Int32
	discComplete atomic.Bool
	finished     atomic.Bool
}

func (n *dirNode) addChild() {
	n.expected.Add(1)
}

func (n *dirNode) childDone() {
	n.done.Add(1)
	n.tryComplete()
}

func (n *dirNode) closeDiscovery() {
	n.discComplete.Store(true)
	n.tryComplete()
}

func (n *dirNode) tryComplete() {
	if !n.discComplete.Load() {
		return
	}
	if n.done.Load() < n.expected.Load() {
		return
	}
	if !n.finished.CompareAndSwap(false, true) {
		return
	}

	if n.t.opt.OnDirComplete != nil {
		n.t.opt.OnDirComplete(&Item{
			Path:     n.rel,
			FullPath: n.path,
			Name:     n.name,
			IsDir:    true,
			Depth:    n.depth,
		})
	}

	if n.parent != nil {
		n.parent.childDone()
		return
	}
	// Root: fire global completion callback before closing done so that any
	// observer waking on Done() sees a fully-settled state (OnDirComplete
	// for root already fired above, OnComplete fires here).
	if n.t.opt.OnComplete != nil {
		n.t.opt.OnComplete()
	}
	close(n.t.done)
}
