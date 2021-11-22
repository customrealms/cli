package template

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	tmpl "text/template"
)

const ManifestFilename = "manifest.json"

type Options struct {
	Name string
}

type Manifest struct {
	Files  map[string]bool   `json:"files"`
	Rename map[string]string `json:"rename"`
}

func readManifest(tmpl Template) (*Manifest, error) {

	// Read the manifest file
	manifestFile, err := tmpl.Open(ManifestFilename)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s in template: %s", ManifestFilename, err)
	}
	defer manifestFile.Close()
	manifestBytes, err := io.ReadAll(manifestFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s in template: %s", ManifestFilename, err)
	}

	// Parse the manifest file as a manifest type
	var manifest Manifest
	if err := json.Unmarshal(manifestBytes, &manifest); err != nil {
		return nil, fmt.Errorf("invalid %s format: %s", ManifestFilename, err)
	}

	// Return the manifest object
	return &manifest, nil

}

func isDirEmpty(dir string) (bool, error) {

	// Check for files in the directory
	entries, err := os.ReadDir(dir)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	// Loop through to find non-hidden files
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), ".") {
			return false, nil
		}
	}

	// The directory is empty
	return true, nil

}

func setupDir(dir string) error {

	// Perform a stat on the directory to check its details
	stat, err := os.Stat(dir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	// If the directory doesn't exist, create it
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(dir, 0777); err != nil {
			return fmt.Errorf("failed to create project directory: %s", err)
		}
		return nil
	}

	// If the directory exists, but is a file
	if !stat.IsDir() {
		return fmt.Errorf("project directory already exists, but is a file")
	}

	// If the project directory is not empty
	empty, err := isDirEmpty(dir)
	if err != nil {
		return err
	}
	if !empty {
		return fmt.Errorf("project directory is not empty")
	}

	// No error if we get here
	return nil

}

func Install(tmpl Template, dir string, options *Options) error {

	// Read the manifest
	manifest, err := readManifest(tmpl)
	if err != nil {
		return err
	}

	// Create the project directory
	if err := setupDir(dir); err != nil {
		return err
	}

	// Populate the project directory using the manifest
	if err := populateDir(tmpl, dir, manifest, options); err != nil {
		return err
	}

	// Return without error
	return nil

}

func populateDir(
	tmpl Template,
	dir string,
	manifest *Manifest,
	options *Options,
) error {

	// Loop through the files in the manifest
	for filename, parse := range manifest.Files {
		if err := populateDirWithFile(tmpl, filename, parse, dir, manifest, options); err != nil {
			return err
		}
	}

	// No errors
	return nil

}

func populateDirWithFile(
	tmpl Template,
	filename string,
	parse bool,
	dir string,
	manifest *Manifest,
	options *Options,
) error {

	// Read the file from the template
	from, err := tmpl.Open(filename)
	if err != nil {
		return err
	}
	defer from.Close()

	// Determine the new name for the file
	to := filename
	if manifest.Rename != nil {
		if renameTo, ok := manifest.Rename[filename]; ok {
			to = renameTo
		}
	}

	// Make sure the directory exists for the file
	toDir := path.Join(dir, path.Dir(to))
	if err := os.MkdirAll(toDir, 0777); err != nil {
		return fmt.Errorf("failed to create %q from template: %s", toDir, err)
	}

	// Copy the file contents
	toFile, err := os.Create(path.Join(dir, to))
	if err != nil {
		return fmt.Errorf("failed to create %q from template: %s", to, err)
	}
	defer toFile.Close()

	// Copy the contents from one to the other
	if parse {
		err = copyAndModifyTemplateFile(toFile, from, options)
	} else {
		_, err = io.Copy(toFile, from)
	}
	if err != nil {
		return err
	}

	// No errors with this file
	return nil

}

func copyAndModifyTemplateFile(to io.Writer, from io.Reader, options *Options) error {

	// Read the entire file to a byte slice
	fromBytes, err := io.ReadAll(from)
	if err != nil {
		return err
	}

	// Create template
	t, err := tmpl.New("").Parse(string(fromBytes))
	if err != nil {
		return err
	}

	// Execute the template into the writer
	return t.Execute(to, options)

}
