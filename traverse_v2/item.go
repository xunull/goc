package traverse_v2

import "io/fs"

// Item is what the user's callback receives. It is built from a DirEntry
// (no extra stat call), so Mode contains only what entry.Type() exposes
// (type bits + permission bits if the OS provided them in readdir).
type Item struct {
	Path     string      // path relative to traverse root, "/" separated; empty for root
	FullPath string      // absolute filesystem path
	Name     string      // basename
	Ext      string      // extension including dot, empty if none
	Mode     fs.FileMode // from DirEntry.Type(); not equivalent to os.Stat's mode
	IsDir    bool
	Depth    int // 0 = root; root's direct children are 1
}
