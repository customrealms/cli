package main

import (
	"os"

	"github.com/customrealms/cli/pkg/build"
	"github.com/customrealms/cli/pkg/project"
)

type BuildCmd struct {
	ProjectDir      string `name:"project" short:"p" usage:"plugin project directory" optional:""`
	ApiVersion      string `name:"mc" usage:"Minecraft version number target" optional:""`
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

	// Create the JAR template to build with
	var jarTemplate build.JarTemplate
	if len(c.TemplateJarFile) > 0 {
		jarTemplate = &build.FileJarTemplate{
			Filename: c.TemplateJarFile,
		}
	} else {
		jarTemplate = &build.GitHubJarTemplate{}
	}

	// Create the project
	crProject := project.New(c.ProjectDir)

	// Create the build action
	buildAction := build.BuildAction{
		Project:     crProject,
		JarTemplate: jarTemplate,
		ApiVersion:  c.ApiVersion,
		OutputFile:  c.OutputFile,
	}
	return buildAction.Run(ctx)
}
