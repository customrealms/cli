package main

import (
	"archive/zip"
	"io"
	"os"
	"strings"
)

const TEMPLATE_JAR = "/Users/conner/Projects/customrealms/bukkit-runtime/target/bukkit-runtime-jar-with-dependencies.jar"
const output = "/Users/conner/Desktop/mcserver/plugins/bukkit-runtime-jar-with-dependencies.jar"

func main() {

	// Open the output file for the final JAR
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a reader for the plugin source code
	pluginCode := strings.NewReader("setInterval(() => console.log('YOOOOOO'), 3000)")

	// Define the plugin.yml details for the plugin
	pluginYml := PluginYml{
		Name:       "MyPlugin",
		ApiVersion: "1.17",
		Version:    "0.0.1",
		Main:       "io.customrealms.MainPlugin",
	}

	// Produce the final JAR file
	if err := CreateFinalJar(
		file,
		TEMPLATE_JAR,
		pluginCode,
		&pluginYml,
	); err != nil {
		panic(err)
	}

}

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
