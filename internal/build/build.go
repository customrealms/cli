package build

import (
	"context"
	_ "embed"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/customrealms/cli/internal/minecraft"
	"github.com/customrealms/cli/internal/project"
)

//go:embed config/webpack.config.js
var webpackConfig string

type BuildAction struct {
	Project          project.Project
	JarTemplate      JarTemplate
	MinecraftVersion minecraft.Version
	OutputFile       string
}

func (a *BuildAction) Run(ctx context.Context) error {
	// Create the temp directory for the code output from Webpack.
	// The code will be put into "bundle.js" in that directory
	webpackOutputDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("cr-build-%d-%d", time.Now().Unix(), rand.Uint32()),
	)
	if err := os.MkdirAll(webpackOutputDir, 0777); err != nil {
		return fmt.Errorf("creating webpack output dir: %w", err)
	}
	defer os.RemoveAll(webpackOutputDir)

	// Write the webpack configuration file temporarily
	webpackConfigFile := filepath.Join(webpackOutputDir, "webpack.config.js")
	if err := os.WriteFile(webpackConfigFile, []byte(webpackConfig), 0777); err != nil {
		return fmt.Errorf("write webpack config file: %w", err)
	}

	// Parse the plugin.yml file
	pluginYML, err := a.Project.PluginYML()
	if err != nil {
		return fmt.Errorf("parse plugin.yml: %w", err)
	}

	fmt.Println("============================================================")
	fmt.Println("Bundling JavaScript code using Webpack")
	fmt.Println("============================================================")

	// Determine the entrypoint for the TypeScript project
	var entrypoint string
	if pluginYML != nil && strings.HasSuffix(pluginYML.Main, ".ts") {
		entrypoint = pluginYML.Main
	} else {
		entrypoint = "./src/main.ts"
	}

	// Build the local directory
	err = a.Project.Exec(ctx, "npx", "webpack-cli",
		"--mode=production",
		"-o", webpackOutputDir,
		"-c", webpackConfigFile,
		"--entry", entrypoint,
	)
	if err != nil {
		return fmt.Errorf("run webpack: %w", err)
	}

	fmt.Println()

	// Package the jar file
	ja := JarAction{
		Project:          a.Project,
		JarTemplate:      a.JarTemplate,
		MinecraftVersion: a.MinecraftVersion,
		BundleFile:       filepath.Join(webpackOutputDir, "bundle.js"),
		OutputFile:       a.OutputFile,
	}
	return ja.Run(ctx)
}
