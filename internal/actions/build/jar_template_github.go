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

func (t *GitHubJarTemplate) getJarUrl() string {
	return fmt.Sprintf(
		"https://github.com/customrealms/bukkit-runtime/releases/latest/download/bukkit-runtime-%s.jar",
		t.MinecraftVersion,
	)
}

func (t *GitHubJarTemplate) Jar() (io.ReadCloser, error) {

	// Get the JAR url
	jarUrl := t.getJarUrl()

	// Download the JAR file
	fmt.Printf(" -> %s\n", jarUrl)
	res, err := http.Get(jarUrl)
	if err != nil {
		return nil, err
	}

	// Return the response body
	return res.Body, nil

}
