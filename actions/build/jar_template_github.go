package build

import (
	"fmt"
	"io"
	"net/http"
	"runtime"

	"github.com/customrealms/cli/minecraft"
)

type GitHubJarTemplate struct {
	OperatingSystem  string
	MinecraftVersion minecraft.Version
}

func (t *GitHubJarTemplate) normalizeOperatingSystem() string {
	if len(t.OperatingSystem) > 0 {
		return t.OperatingSystem
	}
	if runtime.GOOS == "darwin" {
		return "macos"
	}
	return runtime.GOOS
}

func (t *GitHubJarTemplate) getJarUrl() string {
	return fmt.Sprintf(
		"https://github.com/customrealms/bukkit-runtime/releases/latest/download/bukkit-runtime-%s-%s.jar",
		t.normalizeOperatingSystem(),
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
