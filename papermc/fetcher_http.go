package papermc

import (
	"io"
	"net/http"
)

type HttpFetcher struct{}

func (f *HttpFetcher) Fetch(version *Version) (io.ReadCloser, error) {

	// Fetch the JAR file from the server
	res, err := http.Get(version.Url())
	if err != nil {
		return nil, err
	}

	// Return the body of the response
	return res.Body, nil

}
