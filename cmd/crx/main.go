package main

import (
	"fmt"
	"os"
)

const VERSION = "0.1.0"

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
		fmt.Println("crx build ... is not yet implemented")
		os.Exit(1)
	case "build":
		crxBuild()
	}

}
