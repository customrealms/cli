package template

import (
	"io"
	"io/fs"
	"path/filepath"
)

type Template interface {
	Open(name string) (io.ReadCloser, error)
}

type templateFS struct {
	Dir string
	FS  fs.FS
}

func (t *templateFS) Open(name string) (io.ReadCloser, error) {
	return t.FS.Open(filepath.Join(t.Dir, name))
}

func NewFromFS(fileSystem fs.FS) Template {
	return &templateFS{FS: fileSystem}
}

func NewFromFSDir(fileSystem fs.FS, dir string) Template {
	return &templateFS{
		Dir: dir,
		FS:  fileSystem,
	}
}
