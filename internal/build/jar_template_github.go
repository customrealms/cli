package build

import (
	"fmt"
	"io"
	"net/http"

	"github.com/customrealms/cli/internal/minecraft"
)

type GitHubJarTemplate struct {
	MinecraftVersion minecraft.Version
}

func (t *GitHubJarTemplate) Jar() (io.ReadCloser, error) {

	// Get the JAR url
	jarUrl := t.MinecraftVersion.PluginJarUrl()

	// Download the JAR file
	fmt.Printf(" -> %s\n", jarUrl)
	res, err := http.Get(jarUrl)
	if err != nil {
		return nil, err
	}

	// Return the response body
	return res.Body, nil

}
