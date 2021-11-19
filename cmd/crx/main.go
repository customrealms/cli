package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/customrealms/cli/actions/build"
	"github.com/customrealms/cli/actions/initialize"
)

const VERSION = "0.3.0"
const DEFAULT_MC_VERSION = "1.17.1"
const CR_CORE_VERSION_TARGET = "^0.1.0"

func main() {

	// If there are no command line arguments
	if len(os.Args) <= 1 {
		fmt.Println("Missing command name: build, init")
		os.Exit(1)
	}

	// Get the operation string
	var err error
	switch os.Args[1] {
	case "version":
		fmt.Printf("customrealms-cli (crx) v%s\n", VERSION)
	case "init":
		err = crxInit()
	case "build":
		err = crxBuild()
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func crxInit() error {

	// Parse command line arguments
	cwd, _ := os.Getwd()
	var projectDir string
	flag.StringVar(&projectDir, "p", cwd, "plugin project directory")
	flag.CommandLine.Parse(os.Args[2:])

	// Create the init runner
	initAction := initialize.InitAction{
		Name:        filepath.Base(projectDir),
		Dir:         projectDir,
		Template:    nil,
		CoreVersion: CR_CORE_VERSION_TARGET,
		CliVersion:  fmt.Sprintf("^%s", VERSION),
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
	flag.StringVar(&projectDir, "p", ".", "plugin project directory")
	flag.StringVar(&mcVersion, "mc", DEFAULT_MC_VERSION, "Minecraft version number target")
	flag.StringVar(&outputFile, "o", "", "output JAR file path")
	flag.StringVar(&operatingSystem, "os", runtime.GOOS, "operating system target (windows, macos, or linux)")
	flag.CommandLine.Parse(os.Args[2:])

	// Require the output file path
	if len(outputFile) == 0 {
		fmt.Println("Output JAR file is required: -o option")
		os.Exit(1)
	}

	// Create the JAR template to build with
	jarTemplate := build.JarTemplate{
		MinecraftVersion: mcVersion,
		OperatingSystem:  operatingSystem,
	}

	// Create the build action
	buildAction := build.BuildAction{
		ProjectDir:       projectDir,
		JarTemplate:      &jarTemplate,
		MinecraftVersion: mcVersion,
		OutputFile:       outputFile,
	}

	// Run the build action
	return buildAction.Run(context.Background())

}
