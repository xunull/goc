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
	Root         string
	rootTruePath string
	DbName       string
	CacheName    string
	LogName      string
	AssetName    string
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
	m.rootTruePath = h
	return m
}

func (m *MyHome) GetDbHome() string {
	return filepath.Join(m.rootTruePath, m.DbName)
}

func (m *MyHome) GetLogHome() string {
	return filepath.Join(m.rootTruePath, m.LogName)
}

func (m *MyHome) GetAssetHome() string {
	return filepath.Join(m.rootTruePath, m.AssetName)
}

func (m *MyHome) makeHome() {
	makeDir(m.rootTruePath)
}

func (m *MyHome) makeHomeDb() {
	p := filepath.Join(m.rootTruePath, m.DbName)
	makeDir(p)
}

func (m *MyHome) makeHomeCache() {
	p := filepath.Join(m.rootTruePath, m.CacheName)
	makeDir(p)
}

func (m *MyHome) makeLog() {
	p := filepath.Join(m.rootTruePath, m.LogName)
	makeDir(p)
}

func (m *MyHome) makeAsset() {
	p := filepath.Join(m.rootTruePath, m.AssetName)
	makeDir(p)
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
