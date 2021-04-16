package tree

import (
	"context"
	"github.com/xunull/goc/commonx"
	"io/ioutil"
	"path"
	"path/filepath"
	"sync"
)

type FileItem struct {
	Path string
}

type dirTreeWorker struct {
	Option   *Option
	Path     string
	wg       *sync.WaitGroup
	ItemChan chan *FileItem
	errChan  chan error
}

func (w *dirTreeWorker) makeDirPath(ctx context.Context, p, parent string) {
	defer w.wg.Done()

	files, err := ioutil.ReadDir(p)
	if err != nil {
		w.errChan <- err
		return
	}

	for _, file := range files {
		if file.IsDir() {
			w.wg.Add(1)
			go w.makeDirPath(ctx, filepath.Join(p, file.Name()), path.Join(parent, file.Name()))
		} else {
			w.ItemChan <- &FileItem{Path: path.Join(parent, file.Name())}
		}
	}
}

func (w *dirTreeWorker) handle() *TreeElem {

	defer close(w.errChan)

	rootName := filepath.Base(w.Path)
	files, err := ioutil.ReadDir(w.Path)
	commonx.CheckErrOrFatal(err)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		<-w.errChan
		cancel()
	}()

	pathList := make([]interface{}, 0)

	for _, file := range files {
		if file.IsDir() {
			w.wg.Add(1)
			go w.makeDirPath(ctx, filepath.Join(w.Path, file.Name()), rootName)
		} else {
			w.ItemChan <- &FileItem{Path: path.Join(rootName, file.Name())}
		}
	}
	w.wg.Wait()
	close(w.ItemChan)
	for item := range w.ItemChan {
		pathList = append(pathList, item)
	}

	ppt := PlantPathTree(pathList, "Path")
	return ppt

}

type WithOption func(o *Option)

func DirTree(p string, opts ...WithOption) *TreeElem {

	d := GetDefaultOption()
	for _, o := range opts {
		o(d)
	}

	w := dirTreeWorker{Option: d,
		Path:     p,
		ItemChan: make(chan *FileItem, 1024),
		wg:       new(sync.WaitGroup),
		errChan:  make(chan error)}

	return w.handle()

}
