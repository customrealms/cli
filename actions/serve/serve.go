package serve

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type ServeAction struct {
	PaperVersion  *PaperVersion
	PluginJarPath string
}

func (a *ServeAction) DownloadJarTo(dest string) error {
	url := a.PaperVersion.Url()
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()
	if _, err := io.Copy(file, res.Body); err != nil {
		return err
	}
	return nil
}

func copyFile(from, to string) error {
	fromFile, err := os.Open(from)
	if err != nil {
		return err
	}
	defer fromFile.Close()

	toFile, err := os.Create(to)
	if err != nil {
		return err
	}
	defer toFile.Close()

	if _, err := io.Copy(toFile, fromFile); err != nil {
		return err
	}
	return nil
}

func (a *ServeAction) Run(ctx context.Context) error {

	fmt.Println("============================================================")
	fmt.Println("Setting up Minecraft server directory...")
	fmt.Println("============================================================")

	// Create the temp directory
	dir, err := os.MkdirTemp("", "cr-server-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	fmt.Println(" -> ", dir)
	fmt.Println()

	fmt.Println("============================================================")
	fmt.Println("Downloading JAR file for PaperMC server...")
	fmt.Println("============================================================")

	// Create the name of the JAR file
	jarBase := fmt.Sprintf("paper-%s-%d.jar", a.PaperVersion.Version, a.PaperVersion.Build)
	jarFile := filepath.Join(dir, jarBase)

	// Download the JAR file to the path
	if err := a.DownloadJarTo(jarFile); err != nil {
		return err
	}

	fmt.Println(" -> Done")
	fmt.Println()

	fmt.Println("============================================================")
	fmt.Println("Copying plugin JAR file to server 'plugins' folder...")
	fmt.Println("============================================================")

	// Make the plugin directory
	pluginsDir := filepath.Join(dir, "plugins")
	if err := os.MkdirAll(pluginsDir, 0777); err != nil {
		return err
	}
	if err := copyFile(a.PluginJarPath, filepath.Join(pluginsDir, filepath.Base(a.PluginJarPath))); err != nil {
		return err
	}

	// Create the "eula.txt" file
	if err := os.WriteFile(filepath.Join(dir, "eula.txt"), []byte("eula=true\n"), 0777); err != nil {
		return err
	}

	fmt.Println(" -> Done")
	fmt.Println()

	fmt.Println("============================================================")
	fmt.Println("Launching server...")
	fmt.Println("============================================================")
	fmt.Println()

	// Run the server
	cmd := exec.CommandContext(ctx, "java", "-jar", jarBase, "-nogui")
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()

}
