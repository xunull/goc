package dir_tree

import (
	"github.com/xunull/goc/easy/routine_pool"
	"sync"
)

type (
	DTree struct {
		Root        string
		routinePool *routine_pool.RoutinePool
		option      *option
		hf          HandlerFunc
	}

	HandlerFunc struct {
		DirFunc    func(obj *TreeItem)
		FileFunc   func(obj *TreeItem)
		FinishFunc func()
	}
)

func NewTree(root string, hf HandlerFunc, opts ...Option) *DTree {
	dt := &DTree{
		Root: root,
		hf:   hf,
	}
	dt.setOption(opts...)
	dt.routinePool = routine_pool.NewPool(dt.option.WorkerCount)
	return dt
}

func (dt *DTree) setOption(opts ...Option) {
	for _, o := range opts {
		o(dt.option)
	}
}

func (dt *DTree) Exec() {

	dt.routinePool.Start()

	wt := walkTarget{
		dirname: dt.Root,
		dt:      dt,
	}
	var wg sync.WaitGroup
	wg.Add(1)
	wt.pwg = &wg
	wt.walk()
	wg.Wait()

	if dt.hf.FinishFunc != nil {
		dt.hf.FinishFunc()
	}

}
