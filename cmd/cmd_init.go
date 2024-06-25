package main

import (
	"os"
	"path/filepath"

	"github.com/customrealms/cli/internal/initialize"
)

type InitCmd struct {
	ProjectDir string `name:"project" short:"p" usage:"plugin project directory" optional:""`
}

func (c *InitCmd) Run() error {
	// Root context for the CLI
	ctx, cancel := rootContext()
	defer cancel()

	// Default to the current working directory
	if c.ProjectDir == "" {
		c.ProjectDir, _ = os.Getwd()
	}

	// Create the init runner
	initAction := initialize.InitAction{
		Name:     filepath.Base(c.ProjectDir),
		Dir:      c.ProjectDir,
		Template: nil,
	}
	return initAction.Run(ctx)
}
