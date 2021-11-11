package lib

import (
	"archive/zip"
	"io"
	"strings"
)

func CreateFinalJar(
	writer io.Writer,
	templateJar string,
	pluginSourceCode io.Reader,
	pluginYml *PluginYml,
) error {

	// Create the ZIP writer
	zw := zip.NewWriter(writer)
	defer zw.Close()

	// Create the ZIP reader from the base YML
	zr, err := zip.OpenReader(templateJar)
	if err != nil {
		return err
	}

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

	// Write the plugin code to the jar
	codeFile, err := zw.Create("plugin.js")
	if err != nil {
		return err
	}
	if _, err := io.Copy(codeFile, pluginSourceCode); err != nil {
		return err
	}

	// Write the plugin YML file to the jar
	ymlFile, err := zw.Create("plugin.yml")
	if err != nil {
		return err
	}
	if _, err := io.Copy(ymlFile, strings.NewReader(pluginYml.String())); err != nil {
		return err
	}

	// We're done, no errors
	return nil

}
