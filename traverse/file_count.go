package traverse

import (
	"github.com/xunull/goc/lang_ext"
	"sync"
)

type FileCountRes struct {
	Count       int
	TargetCount int
	CountMap    map[string]int
	*option
	mux sync.Mutex
}

func (s *FileCountRes) callbackForGetFileCount(item *TraverseItem) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.Count += 1
	if _, ok := lang_ext.CommonLanguageExt[item.Ext]; ok {
		if _, ok = s.CountMap[lang_ext.CommonLanguageExt[item.Ext]]; ok {
			s.CountMap[lang_ext.CommonLanguageExt[item.Ext]] += 1
		} else {
			s.CountMap[lang_ext.CommonLanguageExt[item.Ext]] = 1
		}
	}
	if s.option.TargetExt != "" {
		if item.Ext == s.option.TargetExt {
			s.TargetCount += 1
		}
	}
}

func GetFileCount(dir string, opts ...Option) *FileCountRes {
	op := &option{}

	for _, o := range opts {
		o(op)
	}

	res := &FileCountRes{
		CountMap: make(map[string]int),
	}
	res.option = op

	t := NewDirTraverse(dir, res.callbackForGetFileCount)
	t.Handle(opts...)
	t.WorkSheet.Wait()
	return res
}
