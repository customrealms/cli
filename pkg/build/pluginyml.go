package build

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/customrealms/cli/pkg/pluginyml"
	"github.com/customrealms/cli/pkg/project"
)

const JarMainClass = "io.customrealms.MainPlugin"

func GeneratePluginYML(project project.Project, apiVersion string, logWriter io.Writer) (*pluginyml.Plugin, error) {
	if logWriter == nil {
		logWriter = os.Stdout
	}

	// Read the package.json file
	packageJSON, err := project.PackageJSON()
	if err != nil {
		return nil, fmt.Errorf("getting package.json: %w", err)
	}

	// Read the plugin.yml file
	plugin, err := project.PluginYML()
	if err != nil {
		return nil, fmt.Errorf("getting plugin.yml: %w", err)
	}

	// If plugin.yml and package.json are both missing, it's an error
	if packageJSON == nil && plugin == nil {
		return nil, errors.New("missing both package.json and plugin.yml")
	}

	// If there is no plugin.yml file present, create one
	if plugin == nil {
		plugin = &pluginyml.Plugin{}
		plugin.Name = packageJSON.Name
	}

	// Set the main Java class for the plugin
	plugin.Main = JarMainClass

	// Set the Bukkit API version for the plugin
	if apiVersion != "" {
		plugin.ApiVersion = &apiVersion
	}

	// If there is a package.json file
	if packageJSON != nil {
		// Update the version if it's missing
		if plugin.Version == "" && packageJSON.Version != "" {
			plugin.Version = packageJSON.Version
		} else if plugin.Version == "" && packageJSON.Version == "" {
			fmt.Fprintln(logWriter, "No version found in plugin.yml or package.json. Consider adding a version to package.json.")
			fmt.Fprintln(logWriter, "Using version '0.0.0' as a fallback.")
			plugin.Version = "0.0.0"
		} else if plugin.Version != packageJSON.Version {
			fmt.Fprintln(logWriter, "Version mismatch between plugin.yml and package.json. Consider removing `version` from plugin.yml.")
			fmt.Fprintf(logWriter, "Using version '%s' from plugin.yml", plugin.Version)
		}
	}

	// Add libraries needed for GraalJS
	plugin.Libraries = append(plugin.Libraries,
		"org.graalvm.polyglot:polyglot:24.1.2",
		"org.graalvm.sdk:collections:24.1.2",
		"org.graalvm.sdk:nativeimage:24.1.2",
		"org.graalvm.sdk:word:24.1.2",
		"org.graalvm.js:js-language:24.1.2",
		"org.graalvm.regex:regex:24.1.2",
		"org.graalvm.truffle:truffle-api:24.1.2",
		"org.graalvm.shadowed:icu4j:24.1.2",
		// "org.graalvm.truffle:truffle-runtime:24.1.2",
		"org.graalvm.sdk:jniutils:24.1.2",
		"org.graalvm.truffle:truffle-compiler:24.1.2",
	)

	// Return the plugin yml
	return plugin, nil
}
