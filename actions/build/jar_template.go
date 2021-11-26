package build

import (
	"io"
)

type JarTemplate interface {
	Jar() (io.ReadCloser, error)
}
