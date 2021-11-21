package build

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const JAR_MAIN_CLASS = "io.customrealms.MainPlugin"

type PackageJSON struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func mcVersionToApiVersion(mcVersion string) string {
	parts := strings.Split(mcVersion, ".")
	return strings.Join(parts[:2], ".")
}

type BuildAction struct {
	ProjectDir       string
	JarTemplate      *JarTemplate
	MinecraftVersion string
	OutputFile       string
}

func (a *BuildAction) Run(ctx context.Context) error {

	// Build the local directory
	cmd := exec.CommandContext(ctx, "npx", "webpack-cli", "--mode=production")
	cmd.Dir = a.ProjectDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Download the jar template
	var jarTemplateBuf bytes.Buffer
	if err := a.JarTemplate.Download(&jarTemplateBuf); err != nil {
		return err
	}

	// Make sure the directory above the output file exists
	if err := os.MkdirAll(filepath.Dir(a.OutputFile), 0777); err != nil {
		return err
	}

	// Open the output file for the final JAR
	file, err := os.Create(a.OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a reader for the plugin source code
	pluginCode, err := os.Open(filepath.Join(a.ProjectDir, "dist", "bundle.js"))
	if err != nil {
		return err
	}
	defer pluginCode.Close()

	// Read the package.json file
	packageJson, err := a.readPackageJson()
	if err != nil {
		return err
	}

	// Define the plugin.yml details for the plugin
	pluginYml := PluginYml{
		Name:       packageJson.Name,
		ApiVersion: mcVersionToApiVersion(a.MinecraftVersion),
		Version:    packageJson.Version,
		Main:       JAR_MAIN_CLASS,
	}

	// Produce the final JAR file
	if err := WriteJarFile(
		file,
		jarTemplateBuf.Bytes(),
		pluginCode,
		&pluginYml,
	); err != nil {
		return err
	}

	fmt.Println("Wrote JAR file to: ", a.OutputFile)

	return nil

}

// readPackageJson reads the package.json file in the project
func (a *BuildAction) readPackageJson() (*PackageJSON, error) {

	// Read the package.json file to bytes
	packageJsonBytes, err := os.ReadFile(filepath.Join(a.ProjectDir, "package.json"))
	if err != nil {
		return nil, err
	}

	// Unmarshal it from its JSON format
	var packageJson PackageJSON
	if err := json.Unmarshal(packageJsonBytes, &packageJson); err != nil {
		return nil, err
	}

	// Return the package JSON object
	return &packageJson, nil

}

func WriteJarFile(
	writer io.Writer,
	templateJarData []byte,
	pluginSourceCode io.Reader,
	pluginYml *PluginYml,
) error {

	fmt.Println("============================================================")
	fmt.Println("Generating final JAR file for your plugin")
	fmt.Println("============================================================")

	// Create the ZIP writer
	zw := zip.NewWriter(writer)
	defer zw.Close()

	// Create the ZIP reader from the base YML
	templateJarReader := bytes.NewReader(templateJarData)
	zr, err := zip.NewReader(templateJarReader, int64(len(templateJarData)))
	if err != nil {
		return err
	}

	fmt.Println(" -> Copying template files to new JAR file")

	// Copy all the files back to the jar file
	for _, f := range zr.File {

		// Skip some files
		if f.Name == "plugin.js" || f.Name == "plugin.yml" {
			continue
		}

		// Copy the rest
		if err := zw.Copy(f); err != nil {
			return err
		}

	}

	fmt.Println(" -> Writing bundle JS code to JAR file")

	// Write the plugin code to the jar
	codeFile, err := zw.Create("plugin.js")
	if err != nil {
		return err
	}
	if _, err := io.Copy(codeFile, pluginSourceCode); err != nil {
		return err
	}

	fmt.Println(" -> Writing plugin.yml file to JAR file")

	// Write the plugin YML file to the jar
	ymlFile, err := zw.Create("plugin.yml")
	if err != nil {
		return err
	}
	if _, err := io.Copy(ymlFile, strings.NewReader(pluginYml.String())); err != nil {
		return err
	}

	fmt.Println(" -> DONE")
	fmt.Println()

	// We're done, no errors
	return nil

}
