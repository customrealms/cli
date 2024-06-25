package server

import (
	"io"

	"github.com/customrealms/cli/internal/minecraft"
)

type JarFetcher interface {
	Fetch(version minecraft.Version) (io.ReadCloser, error)
}
