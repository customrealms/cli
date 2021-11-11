package main

import (
	"os"
	"strings"

	"github.com/customrealms/cli/lib"
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
	pluginYml := lib.PluginYml{
		Name:       "MyPlugin",
		ApiVersion: "1.17",
		Version:    "0.0.1",
		Main:       "io.customrealms.MainPlugin",
	}

	// Produce the final JAR file
	if err := lib.CreateFinalJar(
		file,
		TEMPLATE_JAR,
		pluginCode,
		&pluginYml,
	); err != nil {
		panic(err)
	}

}
