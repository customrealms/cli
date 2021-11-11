package lib

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
)

type JarTemplate struct {
	Platform  string
	MCVersion string
}

func (jt *JarTemplate) normalizePlatform() string {
	if len(jt.Platform) > 0 {
		return jt.Platform
	}
	if runtime.GOOS == "darwin" {
		return "macos"
	}
	return runtime.GOOS
}

func (jt *JarTemplate) getJarUrl() string {
	return fmt.Sprintf(
		"https://github.com/customrealms/bukkit-runtime/releases/latest/download/bukkit-runtime-%s-%s.jar",
		jt.normalizePlatform(),
		jt.MCVersion,
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
