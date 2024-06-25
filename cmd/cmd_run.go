package main

import (
	"os"

	"github.com/customrealms/cli/internal/actions/build"
	"github.com/customrealms/cli/internal/actions/serve"
	"github.com/customrealms/cli/internal/project"
	"github.com/customrealms/cli/internal/server"
)

type RunCmd struct {
	ProjectDir      string `name:"project" short:"p" usage:"plugin project directory" optional:""`
	McVersion       string `name:"mc" short:"mc" usage:"Minecraft version number target" optional:""`
	TemplateJarFile string `name:"jar" short:"t" usage:"template JAR file" optional:""`
}

func (c *RunCmd) Run() error {
	// Root context for the CLI
	ctx, cancel := rootContext()
	defer cancel()

	// Default to the current working directory
	if c.ProjectDir == "" {
		c.ProjectDir, _ = os.Getwd()
	}

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(c.McVersion)

	// Generate a temp filename for the plugin JAR file
	ofile, _ := os.CreateTemp("", "cr-jar-output-*.jar")
	ofile.Close()
	outputFile := ofile.Name()
	defer os.Remove(outputFile)

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
	crProject := project.Project{
		Dir: c.ProjectDir,
	}

	// Create the build action
	buildAction := build.BuildAction{
		Project:          &crProject,
		JarTemplate:      jarTemplate,
		MinecraftVersion: minecraftVersion,
		OutputFile:       outputFile,
	}

	// Run the build action
	if err := buildAction.Run(ctx); err != nil {
		return err
	}

	// Create a fetcher for the Minecraft server JAR file that caches the files locally
	serverJarFetcher, err := server.NewCachedFetcher(&server.HttpFetcher{})
	if err != nil {
		return err
	}

	// Create the serve runner
	serveAction := serve.ServeAction{
		MinecraftVersion: minecraftVersion,
		PluginJarPath:    outputFile,
		ServerJarFetcher: serverJarFetcher,
	}
	return serveAction.Run(ctx)
}
