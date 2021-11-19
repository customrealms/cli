package build

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
)

type JarTemplate struct {
	OperatingSystem  string
	MinecraftVersion string
}

func (jt *JarTemplate) normalizeOperatingSystem() string {
	if len(jt.OperatingSystem) > 0 {
		return jt.OperatingSystem
	}
	if runtime.GOOS == "darwin" {
		return "macos"
	}
	return runtime.GOOS
}

func (jt *JarTemplate) getJarUrl() string {
	return fmt.Sprintf(
		"https://github.com/customrealms/bukkit-runtime/releases/latest/download/bukkit-runtime-%s-%s.jar",
		jt.normalizeOperatingSystem(),
		jt.MinecraftVersion,
	)
}

func (jt *JarTemplate) Download(writer io.Writer) error {

	fmt.Println("============================================================")
	fmt.Println("Downloading JAR plugin runtime")
	fmt.Println("============================================================")

	// Get the JAR url
	jarUrl := jt.getJarUrl()

	// Download the JAR file
	fmt.Printf(" -> %s\n", jarUrl)
	res, err := http.Get(jarUrl)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Pipe the response body to the writer
	if _, err := io.Copy(writer, res.Body); err != nil {
		return err
	}

	fmt.Println(" -> DONE")
	fmt.Println()

	// Return without error
	return nil

}
