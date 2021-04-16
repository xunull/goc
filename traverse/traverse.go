package traverse

import (
	"github.com/xunull/goc/commonx"
	"github.com/xunull/goc/routine_pool"
	"io/fs"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type DirTraverse struct {
	Path      string
	wg        *sync.WaitGroup
	errChan   chan error
	Over      chan struct{}
	WorkSheet *WorkSheet
	option    *option
	Fc        func(item *TraverseItem)

	routinePool *routine_pool.RoutinePool
}

type TraverseItem struct {
	Fs       fs.FileInfo
	IsDir    bool
	FilePath string
	Path     string
	Parent   string
	Ext      string
	Depth    int
}

func NewDirTraverse(p string, fc func(item *TraverseItem)) *DirTraverse {
	d := &DirTraverse{
		Path:        p,
		wg:          &sync.WaitGroup{},
		errChan:     make(chan error),
		Over:        make(chan struct{}, 1),
		WorkSheet:   NewWorkSheet(),
		option:      getDefaultOption(),
		Fc:          fc,
		routinePool: routine_pool.NewPool(1024),
	}

	d.routinePool.Start()

	return d
}

func (t *DirTraverse) processErr() {
	errList := make([]error, 0)
	for err := range t.errChan {
		errList = append(errList, err)
	}
}

func (t *DirTraverse) WaitOver() {
	<-t.Over
	return
}

func (t *DirTraverse) setOption(opts ...Option) {
	for _, o := range opts {
		o(t.option)
	}
}

func (t *DirTraverse) SetOption(opts ...Option) {
	for _, o := range opts {
		o(t.option)
	}
}

func (t *DirTraverse) wrapCallback(item *TraverseItem) {
	t.Fc(item)
	//if !item.IsDir {
	//	t.WorkSheet.TargetOver(item.Path)
	//}
	t.WorkSheet.TargetOver(item.Path)
}

// ---------------------------------------------------------------------------------------------------------------------

func (t *DirTraverse) traverseDir(p string, parent, parentPath string, depth int) {
	defer t.wg.Done()
	if t.option.Depth != 0 && depth > t.option.Depth {
		return
	}

	files, err := ioutil.ReadDir(p)
	if err != nil {
		t.errChan <- err
		return
	}

	for _, file := range files {

		if file.IsDir() {
			if t.option.DefaultExclude || t.option.DotDirExclude {
				if _, ok := CommonExcludeDir[file.Name()]; ok {
					continue
				} else {
					if t.option.DotDirExclude {
						if strings.HasPrefix(file.Name(), ".") {
							continue
						}
					}
				}
			}
			t.wg.Add(1)

			if t.option.SyncMode {
				t.traverseDir(filepath.Join(p, file.Name()),
					file.Name(),
					path.Join(parentPath, file.Name()),
					depth+1)
			} else {

				func(p, parentPath string, depth int, file fs.FileInfo) {
					t.routinePool.Submit(func() {
						t.traverseDir(filepath.Join(p, file.Name()),
							file.Name(),
							path.Join(parentPath, file.Name()),
							depth+1)
					})
				}(p, parentPath, depth, file)

			}

		} else {

			if t.option.OnlyDir {
				continue
			}

			if t.option.TargetExt != "" {
				if path.Ext(file.Name()) != t.option.TargetExt {
					continue
				}
			}
			if t.option.DefaultExclude {
				if _, ok := CommonExcludeFileExt[path.Ext(file.Name())]; ok {
					continue
				}
			}
		}

		ti := &TraverseItem{
			Fs:       file,
			FilePath: filepath.Join(p, file.Name()),
			IsDir:    file.IsDir(),
			Path:     path.Join(parentPath, file.Name()),
			Parent:   parent,
			Ext:      path.Ext(file.Name()),
			Depth:    depth + 1,
		}

		//if !ti.IsDir {
		//	t.WorkSheet.Add(ti.Path)
		//}

		t.WorkSheet.Add(ti.Path)

		if t.option.SyncMode {
			t.wrapCallback(ti)
		} else {
			ti := ti
			go t.wrapCallback(ti)
		}
	}
}

func (t *DirTraverse) Handle(opts ...Option) {
	t.setOption(opts...)
	//if t.option.WithProgressBar {
	//	t.WorkSheet.StartProgressBar()
	//}

	defer close(t.errChan)

	go t.processErr()

	rootName := filepath.Base(t.Path)
	t.wg.Add(1)
	go t.traverseDir(t.Path, "", rootName, 0)

	t.wg.Wait()
	t.WorkSheet.AddOver()
	t.Over <- struct{}{}
}

// ---------------------------------------------------------------------------------------------------------------------

func (t *DirTraverse) getChildrenPath(fp string, parentPath string, ch chan string) {

	defer t.wg.Done()

	files, err := ioutil.ReadDir(fp)
	commonx.CheckErrOrFatal(err)

	for _, file := range files {
		curFp := filepath.Join(fp, file.Name())
		curPp := path.Join(parentPath, file.Name())
		if file.IsDir() {
			if t.option.DefaultExclude || t.option.DotDirExclude {
				if _, ok := CommonExcludeDir[file.Name()]; ok {
					// pass
				} else {
					if t.option.DotDirExclude {
						if !strings.HasPrefix(file.Name(), ".") {
							t.wg.Add(1)
							go t.getChildrenPath(curFp, curPp, ch)
						}
					} else {
						t.wg.Add(1)
						go t.getChildrenPath(curFp, curPp, ch)
					}
				}
			} else {
				t.wg.Add(1)
				go t.getChildrenPath(curFp, curPp, ch)
			}
		} else {
			ch <- curPp
		}
	}
}

func (t *DirTraverse) GetAllPath(opts ...Option) []string {
	t.setOption(opts...)
	defer close(t.errChan)
	go t.processErr()
	t.wg.Add(1)
	over := make(chan int)
	all := make(chan string, 1024)
	res := make([]string, 0)
	go func() {
		for p := range all {
			res = append(res, p)
		}
		over <- 1
	}()

	t.getChildrenPath(t.Path, "", all)
	t.wg.Wait()

	close(all)
	<-over
	return res
}
