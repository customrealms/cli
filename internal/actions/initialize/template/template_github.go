package template

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func githubUrl(org, repo, branch string) string {
	return fmt.Sprintf("https://github.com/%s/%s/archive/refs/heads/%s.zip", org, repo, branch)
}

func downloadUrl(url string) (io.ReadCloser, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return res.Body, nil
}

func NewFromGitHub(org, repo, branch string) (Template, error) {

	// Get the github url
	url := githubUrl(org, repo, branch)

	// Download the zip file
	body, err := downloadUrl(url)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	// Read the full body to bytes
	zipBytes, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	// Read the bytes
	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return nil, err
	}

	// Define the base directory in the zip file
	baseDir := fmt.Sprintf("%s-%s", repo, branch)

	// Return the template with the file system
	return NewFromFSDir(zipReader, baseDir), nil

}
