package build

import (
	"errors"
	"fmt"
	"log"

	"github.com/customrealms/cli/internal/minecraft"
	"github.com/customrealms/cli/internal/pluginyml"
	"github.com/customrealms/cli/internal/project"
)

const JarMainClass = "io.customrealms.MainPlugin"

func GeneratePluginYML(project project.Project, version minecraft.Version) (*pluginyml.Plugin, error) {
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
	if version != nil {
		apiVersion := version.ApiVersion()
		plugin.ApiVersion = &apiVersion
	}

	// If there is a package.json file
	if packageJSON != nil {
		// Update the version if it's missing
		if plugin.Version == "" && packageJSON.Version != "" {
			plugin.Version = packageJSON.Version
		} else if plugin.Version == "" && packageJSON.Version == "" {
			log.Println("No version found in plugin.yml or package.json. Consider adding a version to package.json.")
			log.Println("Using version '0.0.0' as a fallback.")
			plugin.Version = "0.0.0"
		} else if plugin.Version != packageJSON.Version {
			log.Println("Version mismatch between plugin.yml and package.json. Consider removing `version` from plugin.yml.")
			log.Printf("Using version '%s' from plugin.yml", plugin.Version)
		}
	}

	// Return the plugin yml
	return plugin, nil
}
