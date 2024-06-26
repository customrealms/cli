package main

import (
	"fmt"
	"os"

	"github.com/customrealms/cli/internal/build"
	"github.com/customrealms/cli/internal/project"
	"gopkg.in/yaml.v3"
)

type YmlCmd struct {
	ProjectDir string `name:"project" short:"p" usage:"plugin project directory" optional:""`
	McVersion  string `name:"mc" usage:"Minecraft version number target" optional:""`
}

func (c *YmlCmd) Run() error {
	// Root context for the CLI
	ctx, cancel := rootContext()
	defer cancel()

	// Default to the current working directory
	if c.ProjectDir == "" {
		c.ProjectDir, _ = os.Getwd()
	}

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(ctx, c.McVersion)

	// Create the project
	crProject := project.New(c.ProjectDir)

	// Generate the plugin.yml file
	pluginYML, err := build.GeneratePluginYML(crProject, minecraftVersion)
	if err != nil {
		return fmt.Errorf("generating plugin.yml: %w", err)
	}

	// Encode it to stdout
	enc := yaml.NewEncoder(os.Stdout)
	enc.SetIndent(2)
	if err := enc.Encode(pluginYML); err != nil {
		return fmt.Errorf("encoding plugin.yml: %w", err)
	}
	return nil
}
