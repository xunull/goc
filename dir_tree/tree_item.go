package dir_tree

import "io/fs"

type (
	TreeItem struct {
		Fs     fs.FileInfo
		Parent string
		Ext    string
		Depth  int
	}
)

func (t *TreeItem) Name() string {
	return t.Fs.Name()
}

func (t *TreeItem) IsFile() bool {
	return !t.Fs.IsDir()
}

func (t *TreeItem) IsDir() bool {
	return t.Fs.IsDir()
}

func (t *TreeItem) IsSymlink() bool {
	return t.Fs.Mode()&fs.ModeSymlink != 0
}
