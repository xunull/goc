package goc_fs

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"path"
)

type GocWebFile struct {
	io.Seeker
	fs.File
}

func (*GocWebFile) Readdir(count int) ([]fs.FileInfo, error) {
	return nil, nil
}

type GocHttpFs struct {
	Em   embed.FS
	Path string
}

func (g *GocHttpFs) Open(name string) (http.File, error) {

	full := path.Join(g.Path, name)
	file, err := g.Em.Open(full)

	wf := &GocWebFile{
		File: file,
	}
	return wf, err
}
