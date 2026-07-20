package build

import (
	"fmt"
	"io"
	"net/http"
)

type GitHubJarTemplate struct{}

func (t *GitHubJarTemplate) Jar() (io.ReadCloser, error) {
	// Get the JAR url
	jarUrl := "https://github.com/customrealms/bukkit-runtime/releases/latest/download/bukkit-runtime-1.16.1.jar"

	// Download the JAR file
	fmt.Printf(" -> %s\n", jarUrl)
	res, err := http.Get(jarUrl)
	if err != nil {
		return nil, err
	}

	// Return the response body
	return res.Body, nil
}
