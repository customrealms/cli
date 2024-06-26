package main

import (
	"os"

	"github.com/customrealms/cli/internal/build"
	"github.com/customrealms/cli/internal/project"
)

type BuildCmd struct {
	ProjectDir      string `name:"project" short:"p" usage:"plugin project directory" optional:""`
	McVersion       string `name:"mc" usage:"Minecraft version number target" optional:""`
	TemplateJarFile string `name:"jar" short:"t" usage:"template JAR file" optional:""`
	OutputFile      string `name:"output" short:"o" usage:"output JAR file path"`
}

func (c *BuildCmd) Run() error {
	// Root context for the CLI
	ctx, cancel := rootContext()
	defer cancel()

	// Default to the current working directory
	if c.ProjectDir == "" {
		c.ProjectDir, _ = os.Getwd()
	}

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(ctx, c.McVersion)

	// Create the JAR template to build with
	var jarTemplate build.JarTemplate
	if len(c.TemplateJarFile) > 0 {
		jarTemplate = &build.FileJarTemplate{
			Filename: c.TemplateJarFile,
		}
	} else {
		jarTemplate = &build.GitHubJarTemplate{
			MinecraftVersion: minecraftVersion,
		}
	}

	// Create the project
	crProject := project.New(c.ProjectDir)

	// Create the build action
	buildAction := build.BuildAction{
		Project:          crProject,
		JarTemplate:      jarTemplate,
		MinecraftVersion: minecraftVersion,
		OutputFile:       c.OutputFile,
	}
	return buildAction.Run(ctx)
}
