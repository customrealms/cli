package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alecthomas/kong"
	"github.com/customrealms/cli/internal/minecraft"
)

var (
	version = "dev"
	// commit  = "none"
	// date    = "unknown"
)

var cli struct {
	VersionCmd VersionCmd `cmd:"" name:"version" help:"Show the version of the CLI."`
	InitCmd    InitCmd    `cmd:"" name:"init" help:"Initialize a new plugin project."`
	BuildCmd   BuildCmd   `cmd:"" name:"build" help:"Build the plugin JAR file."`
	RunCmd     RunCmd     `cmd:"" name:"run" help:"Build and serve the plugin in a Minecraft server."`
	YmlCmd     YmlCmd     `cmd:"" name:"yml" help:"Generate the plugin.yml file."`
}

func rootContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	return signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
}

func main() {
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}

// mustMinecraftVersion takes a user-supplied Minecraft version string and resolves the corresponding minecraft.Version
// instance. If nothing can be found, it exits the process
func mustMinecraftVersion(ctx context.Context, versionString string) minecraft.Version {
	if len(versionString) == 0 {
		versionString = "1.20.6"
	}
	minecraftVersion, err := minecraft.LookupVersion(ctx, versionString)
	if err != nil {
		fmt.Println("Failed to resolve the Minecraft version: ", err)
		os.Exit(1)
	}
	return minecraftVersion
}
