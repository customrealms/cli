package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/customrealms/cli/actions/initialize"
	"github.com/customrealms/cli/lib"
)

const VERSION = "0.3.0"
const DEFAULT_MC_VERSION = "1.17.1"
const CR_CORE_VERSION = "0.1.0"

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
		CoreVersion: CR_CORE_VERSION,
		CliVersion:  VERSION,
	}

	// Run the init action
	return initAction.Run(context.Background())

}

func crxBuild() error {

	var projectDir string
	var mcVersion string
	var outputFile string

	flag.StringVar(&projectDir, "p", ".", "plugin project directory")
	flag.StringVar(&mcVersion, "mc", DEFAULT_MC_VERSION, "Minecraft version number target")
	flag.StringVar(&outputFile, "o", "", "output JAR file path")

	flag.CommandLine.Parse(os.Args[2:])

	if len(outputFile) == 0 {
		fmt.Println("Output JAR file is required: -o option")
		os.Exit(1)
	}

	// Build the local directory
	cmd := exec.Command("npm", "run", "build")
	cmd.Dir = projectDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Define the specification for the JAR template
	jarTemplate := lib.JarTemplate{
		MCVersion: mcVersion,
	}

	return lib.BundleJar(
		projectDir,
		&jarTemplate,
		mcVersion,
		outputFile,
	)

}
