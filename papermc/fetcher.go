package papermc

import "io"

type Fetcher interface {
	Fetch(version *Version) (io.ReadCloser, error)
}
