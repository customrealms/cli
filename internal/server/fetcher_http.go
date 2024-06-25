package server

import (
	"io"
	"net/http"

	"github.com/customrealms/cli/internal/minecraft"
)

type HttpFetcher struct{}

func (f *HttpFetcher) Fetch(version minecraft.Version) (io.ReadCloser, error) {

	// Fetch the JAR file from the server
	res, err := http.Get(version.ServerJarUrl())
	if err != nil {
		return nil, err
	}

	// Return the body of the response
	return res.Body, nil

}
