package pluginyml

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Plugin struct {
	// Name is the name of your plugin.
	Name string `yaml:"name"`
	// Version is the semantic version of the plugin (e.g. '1.4.1').
	Version string `yaml:"version"`
	// ApiVersion is the version of the Bukkit API your plugin is built against.
	ApiVersion *string `yaml:"api-version,omitempty"`
	// Description is a human friendly description of the functionality your plugin provides.
	Description *string `yaml:"description,omitempty"`
	// Load explicitly states when the plugin should be loaded. if not supplied will default to 'postworld'.
	Load *string `yaml:"load,omitempty"`
	// Author uniquely identifies who developed this plugin.
	Author *string `yaml:"author,omitempty"`
	// Authors allows you to list multiple authors, if it is a collaborative project.
	Authors []string `yaml:"authors,flow,omitempty"`
	// Website is the URL to the plugin's or author's website.
	Website *string `yaml:"website,omitempty"`
	// Main points to the class that extends JavaPlugin.
	Main string `yaml:"main"`
	// Prefix is the name to use when logging to console instead of the plugin's name.
	Prefix *string `yaml:"prefix,omitempty"`
	// SoftDepend is a list of plugins that are required for your plugin to have full functionality.
	SoftDepend []string `yaml:"softdepend,flow,omitempty"`
	// LoadBefore is a list of plugins that should be loaded after your plugin.
	LoadBefore []string `yaml:"loadbefore,flow,omitempty"`
	// Libraries is a list of libraries your plugin needs which can be loaded from Maven Central.
	Libraries []string `yaml:"libraries,omitempty"`
	// Commands is a map of command names to command attributes.
	Commands map[string]Command `yaml:"commands,omitempty"`
	// Permissions is a map of permission names to permission attributes.
	Permissions map[string]Permission `yaml:"permissions,omitempty"`
}

type Command struct {
	// Description is a short description of what the command does.
	Description *string `yaml:"description,omitempty"`
	// Aliases is a list of alternate command names a user may use.
	Aliases []string `yaml:"aliases,flow,omitempty"`
	// Permission is the most basic permission node required to use the command.
	Permission *string `yaml:"permission,omitempty"`
	// PermissionMessage is the message to display to a user when they do not have the required permission.
	PermissionMessage *string `yaml:"permission-message,omitempty"`
	// Usage is a short description of how to use this command.
	Usage *string `yaml:"usage,omitempty"`
}

type Permission struct {
	// Description is a short description of what the permission allows.
	Description *string `yaml:"description,omitempty"`
	// Default is the default value of the permission.
	Default *string `yaml:"default,omitempty"`
	// Children allows you to set children for the permission.
	Children map[string]PermissionChild `yaml:"children,omitempty"`
}

type PermissionChild struct {
	// Bool is non-nil if the child permission is a boolean.
	Bool *bool
	// Permission is non-nil if the child permission is a nested permission.
	Permission *Permission
}

func (p PermissionChild) MarshalYAML() (any, error) {
	if p.Bool != nil {
		return *p.Bool, nil
	}
	if p.Permission != nil {
		return *p.Permission, nil
	}
	return nil, nil
}

func (p *PermissionChild) UnmarshalYAML(node *yaml.Node) error {
	if node.Tag == "!!bool" {
		var b bool
		if err := node.Decode(&b); err != nil {
			return err
		}
		p.Bool = &b
		return nil
	}
	if node.Tag == "!!map" {
		var perm Permission
		if err := node.Decode(&perm); err != nil {
			return err
		}
		p.Permission = &perm
		return nil
	}
	return fmt.Errorf("unsupported type for child permission: %s", node.Tag)
}
