package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/customrealms/cli/internal/build"
	"github.com/customrealms/cli/internal/project"
	"github.com/customrealms/cli/internal/serve"
	"github.com/customrealms/cli/internal/server"
	"github.com/fsnotify/fsnotify"
	"golang.org/x/sync/errgroup"
)

type RunCmd struct {
	ProjectDir      string `name:"project" short:"p" usage:"plugin project directory" optional:""`
	McVersion       string `name:"mc" usage:"Minecraft version number target" optional:""`
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
	minecraftVersion := mustMinecraftVersion(ctx, c.McVersion)

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

	eg, ctx := errgroup.WithContext(ctx)

	chanPluginUpdated := make(chan struct{})
	chanServerStopped := make(chan struct{})
	eg.Go(func() error {
		// Create new watcher.
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return err
		}
		defer watcher.Close()

		// Start listening for events.
		go func() {
			for {
				select {
				case event, ok := <-watcher.Events:
					if !ok {
						return
					}
					if event.Has(fsnotify.Write) {
						// Rebuild the plugin JAR file
						if err := buildAction.Run(ctx); err != nil {
							log.Println("Error: ", err)
						} else {
							select {
							case chanPluginUpdated <- struct{}{}:
							case <-ctx.Done():
								return
							}
						}
					}
				case err, ok := <-watcher.Errors:
					if !ok {
						return
					}
					log.Println("error:", err)
				case <-ctx.Done():
					return
				case <-chanServerStopped:
					return
				}
			}
		}()

		// Add the project directory and src directory to the watcher
		if err := errors.Join(
			watcher.Add(c.ProjectDir),
			watcher.Add(filepath.Join(c.ProjectDir, "src")),
		); err != nil {
			return err
		}

		// Block until the server is stopped
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-chanServerStopped:
			return nil
		}
	})
	eg.Go(func() error {
		defer close(chanServerStopped)

		// Create the serve runner
		serveAction := serve.ServeAction{
			MinecraftVersion: minecraftVersion,
			PluginJarPath:    outputFile,
			ServerJarFetcher: serverJarFetcher,
		}
		return serveAction.Run(ctx, chanPluginUpdated)
	})
	return eg.Wait()
}
