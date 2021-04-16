package file_utils

import (
	"io"
	"os"
)

type LazyWriter struct {
	name string
	f    *os.File
	err  error
}

func (l *LazyWriter) Write(p []byte) (n int, err error) {
	if l.f == nil && l.err == nil {
		l.f, l.err = os.OpenFile(l.name, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	}
	if l.err != nil {
		return 0, l.err
	}
	return l.f.Write(p)
}

func (l *LazyWriter) Close() error {
	if l.f != nil {
		return l.f.Close()
	}
	return nil
}

func NewLazyWriter(name string) io.WriteCloser {
	return &LazyWriter{name: name}
}
