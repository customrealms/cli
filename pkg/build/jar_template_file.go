package build

import (
	"io"
	"os"
)

type FileJarTemplate struct {
	Filename string
}

func (t *FileJarTemplate) Jar() (io.ReadCloser, error) {
	return os.Open(t.Filename)
}
