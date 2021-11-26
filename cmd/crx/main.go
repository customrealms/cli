package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/customrealms/cli/actions/build"
	"github.com/customrealms/cli/actions/initialize"
	"github.com/customrealms/cli/actions/serve"
	"github.com/customrealms/cli/minecraft"
	"github.com/customrealms/cli/project"
	"github.com/customrealms/cli/server"
)

const VERSION = "0.4.7"

func main() {

	// Define the map of commands
	commands := map[string]func() error{
		"version": func() error {
			fmt.Printf("customrealms-cli (crx) v%s\n", VERSION)
			return nil
		},
		"init":  crxInit,
		"serve": crxServe,
		"build": crxBuild,
		"run":   crxBuildAndServe,
		"yml":   crxYaml,
	}

	// If there are no command line arguments
	if len(os.Args) <= 1 {
		fmt.Println("Missing command name.")
		fmt.Println("List of available commands:")
		for cmd := range commands {
			fmt.Printf(" -> crx %s ...\n", cmd)
		}
		fmt.Println()
		os.Exit(1)
	}

	// Check for the command function
	commandFunc, ok := commands[os.Args[1]]
	if !ok {
		fmt.Printf("Unsupported command %q\n", os.Args[1])
		fmt.Println("List of available commands:")
		for cmd := range commands {
			fmt.Printf(" -> crx %s ...\n", cmd)
		}
		fmt.Println()
		os.Exit(1)
	}

	// Run the command function
	if err := commandFunc(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// mustMinecraftVersion takes a user-supplied Minecraft version string and resolves the corresponding minecraft.Version
// instance. If nothing can be found, it exits the process
func mustMinecraftVersion(versionString string) minecraft.Version {
	if len(versionString) == 0 {
		mcVersion := minecraft.LatestVersion()
		if mcVersion == nil {
			fmt.Println("Failed to resolve the default Minecraft version")
			os.Exit(1)
		}
		return mcVersion
	} else {
		minecraftVersion := minecraft.FindVersion(versionString)
		if minecraftVersion == nil {
			fmt.Println("Unsupported Minecraft version: ", versionString)
			fmt.Println()
			fmt.Println("Please use a supported Minecraft version:")
			for _, version := range minecraft.SupportedVersions {
				fmt.Println(" -> ", version)
			}
			fmt.Println()
			os.Exit(1)
		}
		return minecraftVersion
	}
}

func crxYaml() error {

	// Parse command line arguments
	var projectDir string
	var mcVersion string
	flag.StringVar(&projectDir, "p", ".", "plugin project directory")
	flag.StringVar(&mcVersion, "mc", "", "Minecraft version number target")
	flag.CommandLine.Parse(os.Args[2:])

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(mcVersion)

	// Create the project
	crProject := project.Project{
		Dir: projectDir,
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

func crxInit() error {

	// Parse command line arguments
	cwd, _ := os.Getwd()
	var projectDir string
	flag.StringVar(&projectDir, "p", cwd, "plugin project directory")
	flag.CommandLine.Parse(os.Args[2:])

	// Create the init runner
	initAction := initialize.InitAction{
		Name:     filepath.Base(projectDir),
		Dir:      projectDir,
		Template: nil,
	}

	// Run the init action
	return initAction.Run(context.Background())

}

func crxBuild() error {

	// Parse command line arguments
	var projectDir string
	var mcVersion string
	var outputFile string
	var operatingSystem string
	var templateJarFile string
	flag.StringVar(&projectDir, "p", ".", "plugin project directory")
	flag.StringVar(&mcVersion, "mc", "", "Minecraft version number target")
	flag.StringVar(&outputFile, "o", "", "output JAR file path")
	flag.StringVar(&operatingSystem, "os", "", "operating system target (windows, macos, or linux)")
	flag.StringVar(&templateJarFile, "jar", "", "template JAR file")
	flag.CommandLine.Parse(os.Args[2:])

	// Require the output file path
	if len(outputFile) == 0 {
		fmt.Println("Output JAR file is required: -o option")
		os.Exit(1)
	}

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(mcVersion)

	// Create the JAR template to build with
	var jarTemplate build.JarTemplate
	if len(templateJarFile) > 0 {
		jarTemplate = &build.FileJarTemplate{
			Filename: templateJarFile,
		}
	} else {
		jarTemplate = &build.GitHubJarTemplate{
			MinecraftVersion: minecraftVersion,
			OperatingSystem:  operatingSystem,
		}
	}

	// Create the project
	crProject := project.Project{
		Dir: projectDir,
	}

	// Create the build action
	buildAction := build.BuildAction{
		Project:          &crProject,
		JarTemplate:      jarTemplate,
		MinecraftVersion: minecraftVersion,
		OutputFile:       outputFile,
	}

	// Run the build action
	return buildAction.Run(context.Background())

}

func crxServe() error {

	// Parse command line arguments
	var jarFile string
	var mcVersion string
	flag.StringVar(&jarFile, "jar", "", "path to the plugin JAR file")
	flag.StringVar(&mcVersion, "mc", "", "Minecraft version number target")
	flag.CommandLine.Parse(os.Args[2:])

	// Require the JAR file path
	if len(jarFile) == 0 {
		fmt.Println("JAR file path is required: -jar option")
		os.Exit(1)
	}

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(mcVersion)

	// Create the context
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create a fetcher for the Minecraft server JAR file that caches the files locally
	serverJarFetcher, err := server.NewCachedFetcher(&server.HttpFetcher{})
	if err != nil {
		return err
	}

	// Create the serve runner
	serveAction := serve.ServeAction{
		MinecraftVersion: minecraftVersion,
		PluginJarPath:    jarFile,
		ServerJarFetcher: serverJarFetcher,
	}

	// Run the init action
	return serveAction.Run(ctx)

}

func crxBuildAndServe() error {

	// Parse command line arguments
	var projectDir string
	var mcVersion string
	var outputFile string
	var operatingSystem string
	var templateJarFile string
	flag.StringVar(&projectDir, "p", ".", "plugin project directory")
	flag.StringVar(&mcVersion, "mc", "", "Minecraft version number target")
	flag.StringVar(&outputFile, "o", "", "output JAR file path")
	flag.StringVar(&templateJarFile, "jar", "", "template JAR file")
	flag.StringVar(&operatingSystem, "os", "", "operating system target (windows, macos, or linux)")
	flag.CommandLine.Parse(os.Args[2:])

	// Get the Minecraft version
	minecraftVersion := mustMinecraftVersion(mcVersion)

	// If there is no output file provided, default to a temp file
	if len(outputFile) == 0 {

		// Generate a temp filename for the plugin JAR file
		ofile, _ := os.CreateTemp("", "cr-jar-output-*.jar")
		ofile.Close()
		outputFile = ofile.Name()

		// Make sure to delete the generated file at the end
		defer os.Remove(outputFile)

	}

	// Create the context
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Create the JAR template to build with
	var jarTemplate build.JarTemplate
	if len(templateJarFile) > 0 {
		jarTemplate = &build.FileJarTemplate{
			Filename: templateJarFile,
		}
	} else {
		jarTemplate = &build.GitHubJarTemplate{
			MinecraftVersion: minecraftVersion,
			OperatingSystem:  operatingSystem,
		}
	}

	// Create the project
	crProject := project.Project{
		Dir: projectDir,
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

	// Run the init action
	return serveAction.Run(ctx)

}
