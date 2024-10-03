package dir_tree

import (
	"github.com/xunull/goc/easy/routine_pool"
	"io/fs"
)

type (
	DTree struct {
		Root        string
		routinePool *routine_pool.RoutinePool
		option      *option
	}
	TItemInfo struct {
		Fs       fs.FileInfo
		IsDir    bool
		FilePath string
		Path     string
		Parent   string
		Ext      string
		Depth    int
	}
)

func CreateTree(root string, opts ...Option) *DTree {
	dt := &DTree{
		Root: root,
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

}
