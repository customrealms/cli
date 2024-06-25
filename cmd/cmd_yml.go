package main

import (
	"fmt"
	"os"

	"github.com/customrealms/cli/internal/actions/build"
	"github.com/customrealms/cli/internal/project"
)

type YmlCmd struct {
	ProjectDir string `name:"project" short:"p" usage:"plugin project directory" optional:""`
	McVersion  string `name:"mc" short:"mc" usage:"Minecraft version number target" optional:""`
}

func (c *YmlCmd) Run() error {
	// Default to the current working directory
	if c.ProjectDir == "" {
		c.ProjectDir, _ = os.Getwd()
	}

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(c.McVersion)

	// Create the project
	crProject := project.Project{
		Dir: c.ProjectDir,
	}

	// Read the package.json file
	packageJson, err := crProject.PackageJSON()
	if err != nil {
		return err
	}

	// Define the plugin.yml details for the plugin
	pluginYml := &build.PluginYml{
		MinecraftVersion: minecraftVersion,
		PackageJSON:      packageJson,
	}
	fmt.Println(pluginYml)

	return nil
}
