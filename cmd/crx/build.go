package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/customrealms/cli/lib"
)

const DEFAULT_MC_VERSION = "1.17.1"

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
