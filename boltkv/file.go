package boltkv

import (
	"fmt"
	"github.com/xunull/goc/file_path"
	"path/filepath"
)

func getDbFilename(name string) string {
	return fmt.Sprintf("%s.db", name)
}

func getDbFilePath(name string, dir string) string {
	return filepath.Join(dir, getDbFilename(name))
}

func IsDbExist(name string, dir string) (bool, error) {
	return file_path.PathExists(getDbFilePath(name, dir))
}
