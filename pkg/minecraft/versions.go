package minecraft

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type paperMcBuild struct {
	ID        int    `json:"id"`
	Channel   string `json:"channel"`
	Downloads map[string]struct {
		Name      string `json:"name"`
		Checksums struct {
			Sha256 string `json:"sha256"`
		} `json:"checksums"`
		Size int64  `json:"size"`
		URL  string `json:"url"`
	} `json:"downloads"`
}

func LookupVersion(ctx context.Context, versionStr string) (Version, error) {
	// Lookup the version from PaperMC
	builds, err := downloadJSON[[]paperMcBuild](ctx, fmt.Sprintf("https://fill.papermc.io/v3/projects/paper/versions/%s/builds", versionStr))
	if err != nil {
		return nil, fmt.Errorf("download builds list: %w", err)
	}
	if builds == nil || len(*builds) == 0 {
		return nil, fmt.Errorf("no builds found for version %s", versionStr)
	}

	// The latest entry is the latest build
	build := (*builds)[0]

	// Get the server jar URL for the build
	serverJarDownload, ok := build.Downloads["server:default"]
	if !ok {
		return nil, fmt.Errorf("no server jar found for build %d", build.ID)
	}

	version := &paperMcVersion{versionStr, build.ID, serverJarDownload.URL}
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
