package myhome

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/xunull/goc/commonx"
	"github.com/xunull/goc/file_path"
	"os"
	"path/filepath"
)

type MyHome struct {
	Root      string
	DbName    string
	CacheName string
	LogName   string
}

func NewMyHome(root string) *MyHome {
	return &MyHome{
		Root:      root,
		DbName:    "db",
		CacheName: "cache",
		LogName:   "log",
	}
}

func (m *MyHome) GetDbHome() string {
	if filepath.IsAbs(m.Root) {
		return filepath.Join(m.Root, m.DbName)
	} else {
		return filepath.Join(getHomeDir(), m.Root, m.DbName)
	}
}

func (m *MyHome) GetLogHome() string {
	if filepath.IsAbs(m.Root) {
		return filepath.Join(m.Root, m.LogName)
	} else {
		return filepath.Join(getHomeDir(), m.Root, m.LogName)
	}

}

func (m *MyHome) makeHome() {
	if filepath.IsAbs(m.Root) {
		makeDir(m.Root)
	} else {
		p := filepath.Join(getHomeDir(), m.Root)
		makeDir(p)
	}
}

func (m *MyHome) makeHomeDb() {
	if filepath.IsAbs(m.Root) {
		p := filepath.Join(m.Root, m.DbName)
		makeDir(p)
	} else {
		p := filepath.Join(getHomeDir(), m.Root, m.DbName)
		makeDir(p)
	}

}

func (m *MyHome) makeHomeCache() {
	if filepath.IsAbs(m.Root) {
		p := filepath.Join(m.Root, m.CacheName)
		makeDir(p)
	} else {
		p := filepath.Join(getHomeDir(), m.Root, m.CacheName)
		makeDir(p)
	}
}

func (m *MyHome) makeLog() {
	if filepath.IsAbs(m.Root) {
		p := filepath.Join(m.Root, m.LogName)
		makeDir(p)
	} else {
		p := filepath.Join(getHomeDir(), m.Root, m.LogName)
		makeDir(p)
	}
}

func (m *MyHome) Init() {
	m.makeHome()
	m.makeHomeDb()
	m.makeHomeCache()
	m.makeLog()
}

func makeDir(p string) {
	exist, err := file_path.PathExists(p)
	commonx.CheckErrOrFatal(err)
	if !exist {
		err = os.Mkdir(p, 0700)
		commonx.CheckErrOrFatal(err)
	}
}

func getHomeDir() string {
	home, err := homedir.Dir()
	cobra.CheckErr(err)
	return home
}
