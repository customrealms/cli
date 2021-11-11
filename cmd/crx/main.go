package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/customrealms/cli/lib"
)

const VERSION = "0.1.0"
const DEFAULT_MC_VERSION = "1.17.1"

func main() {

	// If there are no command line arguments
	if len(os.Args) <= 1 {
		fmt.Println("Missing command name: build, init")
		os.Exit(1)
	}

	// Get the operation string
	switch os.Args[1] {
	case "version":
		fmt.Printf("customrealms-cli (crx) v%s\n", VERSION)
	case "init":
		crxInit()
	case "build":
		crxBuild()
	}

}

func crxInit() {

	var projectDir string

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.StringVar(&projectDir, "p", cwd, "plugin project directory")

	flag.CommandLine.Parse(os.Args[2:])

	if err := lib.InitDir(projectDir, filepath.Base(projectDir)); err != nil {
		panic(err)
	}

}

func crxBuild() {

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
		panic(err)
	}

	// Define the specification for the JAR template
	jarTemplate := lib.JarTemplate{
		MCVersion: mcVersion,
	}

	if err := lib.BundleJar(
		projectDir,
		&jarTemplate,
		mcVersion,
		outputFile,
	); err != nil {
		panic(err)
	}

}
