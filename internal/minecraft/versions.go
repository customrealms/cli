package minecraft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type paperMcBuilds struct {
	Builds []paperMcBuild `json:"builds"`
}

type paperMcBuild struct {
	Build     int    `json:"build"`
	Channel   string `json:"channel"`
	Downloads struct {
		Application struct {
			Name   string `json:"name"`
			Sha256 string `json:"sha256"`
		} `json:"application"`
	} `json:"downloads"`
}

func LookupVersion(ctx context.Context, versionStr string) (Version, error) {
	// Lookup the version from PaperMC
	builds, err := downloadJSON[paperMcBuilds](ctx, fmt.Sprintf("https://papermc.io/api/v2/projects/paper/versions/%s/builds", versionStr))
	if err != nil {
		return nil, fmt.Errorf("download builds list: %w", err)
	}
	if builds == nil || len(builds.Builds) == 0 {
		return nil, fmt.Errorf("no builds found for version %s", versionStr)
	}

	// The last entry is the latest build
	build := builds.Builds[len(builds.Builds)-1]
	version := &paperMcVersion{versionStr, build.Build}

	// Check that the version has a downloadable plugin JAR
	ok, err := checkHttpOK(ctx, version.PluginJarUrl())
	if err != nil {
		return nil, fmt.Errorf("check plugin jar url: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("no customrealms bukkit-runtime found for version %s", versionStr)
	}
	return version, nil
}

func downloadJSON[T any](ctx context.Context, url string) (*T, error) {
	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}

	// Send the HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}
	defer res.Body.Close()

	// Decode the JSON response
	var result T
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode json response: %w", err)
	}
	return &result, nil
}

func checkHttpOK(ctx context.Context, url string) (bool, error) {
	// Create the HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodHead, url, nil)
	if err != nil {
		return false, fmt.Errorf("create http request: %w", err)
	}

	// Send the HTTP request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, fmt.Errorf("send http request: %w", err)
	}
	defer res.Body.Close()

	// Check the status code
	return res.StatusCode == http.StatusOK, nil
}
