package traverse

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/xunull/goc/easy/routine_pool"
	"github.com/xunull/goc/file_path"
	"github.com/xunull/goc/file_utils"
	"github.com/xunull/goc/lang_ext"
)

type DirTraverse struct {
	Path        string
	wg          *sync.WaitGroup
	errChan     chan error
	Over        chan struct{}
	WorkSheet   *WorkSheet
	option      *option
	ProcessFunc func(item *TraverseItem) // 重命名：Fc -> ProcessFunc
	routinePool *routine_pool.RoutinePool
	ctx         context.Context
	cancel      context.CancelFunc
	errors      []error
	errorsMutex sync.Mutex
}

type TraverseItem struct {
	FileInfo fs.FileInfo // 重命名：Fs -> FileInfo
	IsDir    bool
	FilePath string
	Path     string
	Parent   string
	Ext      string
	Depth    int
}

func NewDirTraverse(p string, processFunc func(item *TraverseItem), opts ...Option) *DirTraverse {
	ctx, cancel := context.WithCancel(context.Background())
	d := &DirTraverse{
		Path:        p,
		wg:          &sync.WaitGroup{},
		errChan:     make(chan error, 100), // 添加缓冲区避免阻塞
		Over:        make(chan struct{}, 1),
		WorkSheet:   NewWorkSheet(),
		option:      getDefaultOption(),
		ProcessFunc: processFunc,
		ctx:         ctx,
		cancel:      cancel,
		errors:      make([]error, 0),
	}

	d.setOption(opts...)

	d.routinePool = routine_pool.NewPool(d.option.WorkerCount)
	d.routinePool.Start()

	return d
}

func (t *DirTraverse) processErr() {
	for err := range t.errChan {
		if err != nil {
			t.errorsMutex.Lock()
			t.errors = append(t.errors, err)
			t.errorsMutex.Unlock()
		}
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
	t.setOption(opts...)
}

func (t *DirTraverse) wrapCallback(item *TraverseItem) {
	t.ProcessFunc(item)
	t.WorkSheet.ItemOver(item.Path)
}

// ---------------------------------------------------------------------------------------------------------------------

// main method
func (t *DirTraverse) traverseDir(p string, parent, parentPath string, depth int) {
	defer t.wg.Done()

	// 检查 context 是否已取消
	select {
	case <-t.ctx.Done():
		return
	default:
	}

	if t.option.Depth != 0 && depth > t.option.Depth {
		return
	}

	entries, err := os.ReadDir(p)
	if err != nil {
		select {
		case t.errChan <- err:
		default:
			// channel 已满，记录到错误列表
			t.errorsMutex.Lock()
			t.errors = append(t.errors, err)
			t.errorsMutex.Unlock()
		}
		return
	}

	// 转换 DirEntry 为 FileInfo
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, info)
	}

	for _, file := range files {

		if file_utils.IsSymlink(file.Mode()) {
			// symlink
			continue
		} else if file.IsDir() {
			if t.option.DefaultExclude || t.option.DotDirExclude {
				if _, ok := lang_ext.CommonExcludeDir[file.Name()]; ok {
					continue
				} else {
					if t.option.DotDirExclude {
						if strings.HasPrefix(file.Name(), ".") {
							continue
						}
					}
				}
			}

			if t.option.ExcludeDir != nil && len(t.option.ExcludeDir) > 0 {
				p := filepath.Join(parentPath, file.Name())
				if _, ok := t.option.excludeDirMap[file_path.RemovePrefixN(p, 1)]; ok {
					continue
				}
			}

			t.wg.Add(1)

			if t.option.SyncMode {
				t.traverseDir(filepath.Join(p, file.Name()),
					file.Name(),
					filepath.Join(parentPath, file.Name()),
					depth+1)
			} else {

				func(p, parentPath string, depth int, file fs.FileInfo) {
					t.routinePool.Submit(func() {
						t.traverseDir(filepath.Join(p, file.Name()),
							file.Name(),
							filepath.Join(parentPath, file.Name()),
							depth+1)
					})
				}(p, parentPath, depth, file)

			}

		} else {

			if t.option.OnlyDir {
				continue
			}
			if t.option.TargetExt != "" {
				if filepath.Ext(file.Name()) != t.option.TargetExt {
					continue
				}
			}
			if t.option.DefaultExclude {
				if _, ok := lang_ext.CommonExcludeFileExt[filepath.Ext(file.Name())]; ok {
					continue
				}
			}
			if t.option.ExcludeSuffixes != nil {
				flag := false
				for _, suffix := range t.option.ExcludeSuffixes {
					if strings.HasSuffix(file.Name(), suffix) {
						flag = true
						break
					}
				}
				if flag {
					continue
				}
			}
			if t.option.ExcludePrefixes != nil {
				flag := false
				for _, prefix := range t.option.ExcludePrefixes {

					if strings.HasPrefix(file.Name(), prefix) {
						flag = true
						break
					}
				}
				if flag {
					continue
				}
			}
		}

		ti := &TraverseItem{
			FileInfo: file,
			FilePath: filepath.Join(p, file.Name()),
			IsDir:    file.IsDir(),
			Path:     filepath.Join(parentPath, file.Name()),
			Parent:   parent,
			Ext:      filepath.Ext(file.Name()),
			Depth:    depth + 1,
		}

		t.WorkSheet.ItemAdd(ti.Path)

		if t.option.SyncFileOpMode {
			t.wrapCallback(ti)
		} else {
			ti := ti
			go t.wrapCallback(ti)
		}
	}
}

func (t *DirTraverse) Handle(opts ...Option) error {
	t.setOption(opts...)
	defer close(t.errChan)

	go t.processErr()

	// last dir name
	root := filepath.Base(t.Path)

	t.wg.Add(1)
	go t.traverseDir(t.Path, "", root, 0)
	t.wg.Wait()

	t.WorkSheet.TraverseOver()
	t.Over <- struct{}{}

	// 返回收集的错误
	t.errorsMutex.Lock()
	defer t.errorsMutex.Unlock()
	if len(t.errors) > 0 {
		return fmt.Errorf("traverse completed with %d errors: first error: %w", len(t.errors), t.errors[0])
	}
	return nil
}

// ---------------------------------------------------------------------------------------------------------------------

func (t *DirTraverse) getChildrenPath(fp string, parentPath string, ch chan string) {

	defer t.wg.Done()

	// 检查 context 是否已取消
	select {
	case <-t.ctx.Done():
		return
	default:
	}

	entries, err := os.ReadDir(fp)
	if err != nil {
		select {
		case t.errChan <- err:
		default:
			t.errorsMutex.Lock()
			t.errors = append(t.errors, err)
			t.errorsMutex.Unlock()
		}
		return
	}

	// 转换 DirEntry 为 FileInfo
	files := make([]fs.FileInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		files = append(files, info)
	}

	for _, file := range files {
		curFp := filepath.Join(fp, file.Name())
		curPp := filepath.Join(parentPath, file.Name())
		if file.IsDir() {
			if t.option.DefaultExclude || t.option.DotDirExclude {
				if _, ok := lang_ext.CommonExcludeDir[file.Name()]; ok {
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

func (t *DirTraverse) GetAllPath(opts ...Option) ([]string, error) {
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

	t.errorsMutex.Lock()
	defer t.errorsMutex.Unlock()
	if len(t.errors) > 0 {
		return res, fmt.Errorf("traverse completed with %d errors", len(t.errors))
	}
	return res, nil
}

// ---------------------------------------------------------------------------------------------------------------------
// 资源管理方法

// Close 清理所有资源
func (t *DirTraverse) Close() error {
	if t.cancel != nil {
		t.cancel()
	}
	if t.routinePool != nil {
		t.routinePool.Release()
	}
	return nil
}

// Cancel 取消遍历操作
func (t *DirTraverse) Cancel() {
	if t.cancel != nil {
		t.cancel()
	}
}

// Errors 返回遍历过程中收集的所有错误
func (t *DirTraverse) Errors() []error {
	t.errorsMutex.Lock()
	defer t.errorsMutex.Unlock()
	// 返回错误副本以避免并发修改
	errorsCopy := make([]error, len(t.errors))
	copy(errorsCopy, t.errors)
	return errorsCopy
}

// HasErrors 检查是否有错误
func (t *DirTraverse) HasErrors() bool {
	t.errorsMutex.Lock()
	defer t.errorsMutex.Unlock()
	return len(t.errors) > 0
}
