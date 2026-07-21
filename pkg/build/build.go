package build

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/customrealms/cli/pkg/project"
	"github.com/evanw/esbuild/pkg/api"
)

type BuildAction struct {
	Project     project.Project
	JarTemplate JarTemplate
	ApiVersion  string
	OutputFile  string
}

func (a *BuildAction) Run(ctx context.Context) error {
	// Create a temporary directory for the output bundle.
	// The code will be written to "bundle.js" in that directory.
	buildOutputDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("cr-build-%d-%d", time.Now().Unix(), rand.Uint32()),
	)
	if err := os.MkdirAll(buildOutputDir, 0777); err != nil {
		return fmt.Errorf("creating build output dir: %w", err)
	}
	defer os.RemoveAll(buildOutputDir)

	// Parse the plugin.yml file
	pluginYML, err := a.Project.PluginYML()
	if err != nil {
		return fmt.Errorf("parse plugin.yml: %w", err)
	}

	fmt.Println("============================================================")
	fmt.Println("Bundling JavaScript code using esbuild")
	fmt.Println("============================================================")

	// Determine the entrypoint for the TypeScript project
	var entrypoint string
	if pluginYML != nil && strings.HasSuffix(pluginYML.Main, ".ts") {
		entrypoint = pluginYML.Main
	} else {
		entrypoint = "./src/main.ts"
	}

	// Build the local directory using esbuild's Go API.
	result := api.Build(api.BuildOptions{
		AbsWorkingDir:     a.Project.Dir(),
		EntryPoints:       []string{entrypoint},
		Outfile:           filepath.Join(buildOutputDir, "bundle.js"),
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		TreeShaking:       api.TreeShakingTrue,
		Platform:          api.PlatformBrowser,
		Format:            api.FormatIIFE,
		Target:            api.ES2015,
		LogLevel:          api.LogLevelInfo,
		Write:             true,
	})
	if len(result.Errors) > 0 {
		return fmt.Errorf("bundle code with esbuild: %s", result.Errors[0].Text)
	}

	fmt.Println()

	// Package the jar file
	ja := JarAction{
		Project:     a.Project,
		JarTemplate: a.JarTemplate,
		BundleFile:  filepath.Join(buildOutputDir, "bundle.js"),
		OutputFile:  a.OutputFile,
	}
	return ja.Run(ctx)
}
