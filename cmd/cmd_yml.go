package main

import (
	"fmt"
	"os"

	"github.com/customrealms/cli/pkg/build"
	"github.com/customrealms/cli/pkg/project"
	"gopkg.in/yaml.v3"
)

type YmlCmd struct {
	ProjectDir string `name:"project" short:"p" usage:"plugin project directory" optional:""`
	ApiVersion string `name:"mc" usage:"Minecraft version number target" optional:""`
}

func (c *YmlCmd) Run() error {
	// Default to the current working directory
	if c.ProjectDir == "" {
		c.ProjectDir, _ = os.Getwd()
	}

	// Create the project
	crProject := project.New(c.ProjectDir)

	// Generate the plugin.yml file
	pluginYML, err := build.GeneratePluginYML(crProject, c.ApiVersion)
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
