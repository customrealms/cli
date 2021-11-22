package papermc

import (
	"io"

	"github.com/customrealms/cli/minecraft"
)

type Fetcher interface {
	Fetch(version minecraft.Version) (io.ReadCloser, error)
}
