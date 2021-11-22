package project

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
)

type Project struct {
	Dir string
}

// CommandContext creates a command that will execute with default settings and within the project directory
func (p *Project) CommandContext(ctx context.Context, name string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = p.Dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// PackageJSON reads the package.json file contents from the project directory
func (p *Project) PackageJSON() (*PackageJSON, error) {

	// Read the package.json file to bytes
	packageJsonBytes, err := os.ReadFile(filepath.Join(p.Dir, "package.json"))
	if err != nil {
		return nil, err
	}

	// Unmarshal it from its JSON format
	var packageJson PackageJSON
	if err := json.Unmarshal(packageJsonBytes, &packageJson); err != nil {
		return nil, err
	}

	// Return the package JSON object
	return &packageJson, nil

}
