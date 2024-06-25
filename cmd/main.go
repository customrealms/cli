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

var cli struct {
	VersionCmd VersionCmd `cmd:"" help:"Show the version of the CLI."`
	InitCmd    InitCmd    `cmd:"" help:"Initialize a new plugin project."`
	BuildCmd   BuildCmd   `cmd:"" help:"Build the plugin JAR file."`
	RunCmd     RunCmd     `cmd:"" help:"Build and serve the plugin in a Minecraft server."`
	YmlCmd     YmlCmd     `cmd:"" help:"Generate the plugin.yml file."`
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
