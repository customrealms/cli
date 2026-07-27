package build

import (
	"io"
	"net/http"
)

type GitHubJarTemplate struct{}

func (t *GitHubJarTemplate) URL() string {
	return "https://github.com/customrealms/bukkit-runtime/releases/latest/download/bukkit-runtime.jar"
}

func (t *GitHubJarTemplate) Jar() (io.ReadCloser, error) {
	// Download the JAR file
	res, err := http.Get(t.URL())
	if err != nil {
		return nil, err
	}

	// Return the response body
	return res.Body, nil
}
