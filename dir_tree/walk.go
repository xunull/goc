package dir_tree

import (
	"github.com/xunull/goc/file_path"
	"github.com/xunull/goc/file_utils"
	"github.com/xunull/goc/lang_ext"
	"os"
	"path"
	"strings"
	"sync"
)

type (
	walkTarget struct {
		dirname string
		dt      *DTree
		wg      *sync.WaitGroup
		pwg     *sync.WaitGroup
		depth   int
	}
)

func (wt *walkTarget) walk() {
	defer wt.pwg.Done()

	if wt.depth > wt.dt.option.Depth {
		return
	}
	entries, err := os.ReadDir(wt.dirname)
	if err != nil {

	}
	dirList := make([]os.DirEntry, 0)
	fileList := make([]os.DirEntry, 0)
	for _, entry := range entries {
		if file_utils.IsSymlink(entry.Type()) {
			continue
		}
		if entry.IsDir() {
			dirList = append(dirList, entry)
		} else {
			fileList = append(fileList, entry)
		}
	}

	if !wt.dt.option.OnlyDir {
		// handle file

		wt.wg.Add(len(fileList))
		// 1.22以上 才可以这样使用
		for _, entry := range fileList {
			wt.dt.routinePool.Submit(func() {
				wt.handleFile(entry)
			})
		}
	}

	if len(dirList) == 0 {
		return
	}

	dirList = wt.filterDir(dirList)
	if len(dirList) == 0 {
		return
	}

	wt.wg.Add(len(dirList))
	for _, entry := range dirList {

		wt.dt.routinePool.Submit(func() {

			if wt.dt.hf.DirFunc != nil {
				wt.dt.hf.DirFunc(wt.createTreeItem(entry))
			}

			wt.createSubWalkTarget(entry.Name()).walk()
		})
	}

	defer wt.wg.Wait()

}

func (wt *walkTarget) createSubWalkTarget(sub string) *walkTarget {

	dirname := path.Join(wt.dirname, sub)
	return &walkTarget{
		dirname: dirname,
		dt:      wt.dt,
		wg:      &sync.WaitGroup{},
		pwg:     wt.wg,
		depth:   wt.depth + 1,
	}

}

func (wt *walkTarget) filterDir(dirList []os.DirEntry) []os.DirEntry {

	res := make([]os.DirEntry, 0)

	for _, entry := range dirList {
		if wt.dt.option.DefaultExclude || wt.dt.option.DotDirExclude {
			if _, ok := lang_ext.CommonExcludeDir[entry.Name()]; ok {
				continue
			} else {
				if wt.dt.option.DotDirExclude {
					if strings.HasPrefix(entry.Name(), ".") {
						continue
					}
				}
			}
		}

		if wt.dt.option.ExcludeDir != nil && len(wt.dt.option.ExcludeDir) > 0 {
			p := path.Join(wt.dirname, entry.Name())
			if _, ok := wt.dt.option.excludeDirMap[file_path.RemovePrefixN(p, 1)]; ok {
				continue
			}
		}

		res = append(res, entry)
	}
	return res

}

func (wt *walkTarget) createTreeItem(entry os.DirEntry) *TreeItem {
	i, err := entry.Info()
	if err != nil {
		// todo
	}
	ti := &TreeItem{
		Fs: i,
	}
	ti.Parent = wt.dirname
	ti.Ext = path.Ext(entry.Name())
	ti.Depth = wt.depth

	return ti
}

func (wt *walkTarget) handleFile(entry os.DirEntry) {

	defer wt.wg.Done()

	// only handle target ext file
	if wt.dt.option.TargetExt != "" {
		if path.Ext(entry.Name()) != wt.dt.option.TargetExt {
			return
		}
	}
	if wt.dt.option.DefaultExclude {
		if _, ok := lang_ext.CommonExcludeFileExt[path.Ext(entry.Name())]; ok {
			return
		}
	}
	// exclude suffixes
	if wt.dt.option.ExcludeSuffixes != nil {
		flag := false
		for _, suffix := range wt.dt.option.ExcludeSuffixes {
			if strings.HasSuffix(entry.Name(), suffix) {
				flag = true
				break
			}
		}
		if flag {
			return
		}
	}
	// exclude prefixes
	if wt.dt.option.ExcludePrefixes != nil {
		flag := false
		for _, prefix := range wt.dt.option.ExcludePrefixes {

			if strings.HasPrefix(entry.Name(), prefix) {
				flag = true
				break
			}
		}
		if flag {
			return
		}
	}
	if wt.dt.hf.FileFunc != nil {
		wt.dt.hf.FileFunc(wt.createTreeItem(entry))
	}

}
