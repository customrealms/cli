package project

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/customrealms/cli/internal/pluginyml"
	"gopkg.in/yaml.v3"
)

type Project interface {
	// Exec executes a command in the project directory.
	Exec(ctx context.Context, name string, args ...string) error
	// PackageJSON reads the package.json file contents from the project directory.
	// If the file does not exist, it returns nil.
	PackageJSON() (*PackageJSON, error)
	// PluginYML reads the plugin.yml file contents from the project directory.
	// If the file does not exist, it returns nil.
	PluginYML() (*pluginyml.Plugin, error)
}

// New creates a new project from the given directory.
func New(dir string) Project {
	return &project{dir}
}

type project struct {
	dir string
}

func (p *project) Exec(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = p.dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (p *project) PackageJSON() (*PackageJSON, error) {
	// Open the file
	file, err := os.Open(filepath.Join(p.dir, "package.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("opening package.json: %w", err)
	}
	defer file.Close()

	// Decode the json file
	var packageJSON PackageJSON
	if err := json.NewDecoder(file).Decode(&packageJSON); err != nil {
		return nil, fmt.Errorf("decoding package.json: %w", err)
	}
	return &packageJSON, nil
}

func (p *project) PluginYML() (*pluginyml.Plugin, error) {
	// Open the file
	file, err := os.Open(filepath.Join(p.dir, "plugin.yml"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("opening plugin.yml: %w", err)
	}
	defer file.Close()

	// Decode the yaml file
	var plugin pluginyml.Plugin
	if err := yaml.NewDecoder(file).Decode(&plugin); err != nil {
		return nil, fmt.Errorf("decoding plugin.yml: %w", err)
	}
	return &plugin, nil
}
