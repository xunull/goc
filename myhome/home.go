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
	Root        string
	rootAbsPath string
	DbName      string
	CacheName   string
	LogName     string
	AssetName   string
}

func NewMyHome(root string) *MyHome {
	m := &MyHome{
		Root:      root,
		DbName:    "db",
		CacheName: "cache",
		LogName:   "log",
		AssetName: "asset",
	}
	h, err := homedir.Expand(m.Root)
	commonx.CheckErrOrFatal(err)
	m.rootAbsPath = h
	return m
}

func (m *MyHome) GetTargetPath(target string) string {
	return filepath.Join(m.rootAbsPath, target)
}

func (m *MyHome) GetDbHome() string {
	return m.GetTargetPath(m.DbName)
}

func (m *MyHome) GetLogHome() string {
	return m.GetTargetPath(m.LogName)
}

func (m *MyHome) GetAssetHome() string {
	return m.GetTargetPath(m.AssetName)
}

func (m *MyHome) makeHome() {
	makeDir(m.rootAbsPath)
}

func (m *MyHome) makeHomeDb() {
	makeDir(m.GetDbHome())
}

func (m *MyHome) makeHomeCache() {
	makeDir(m.GetTargetPath(m.CacheName))
}

func (m *MyHome) makeLog() {
	makeDir(m.GetLogHome())
}

func (m *MyHome) makeAsset() {
	makeDir(m.GetAssetHome())
}

func (m *MyHome) Init() {
	m.makeHome()
	m.makeHomeDb()
	m.makeHomeCache()
	m.makeLog()
	m.makeAsset()
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
