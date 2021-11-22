package build

import (
	"fmt"
	"strings"

	"github.com/customrealms/cli/minecraft"
	"github.com/customrealms/cli/project"
)

const JAR_MAIN_CLASS = "io.customrealms.MainPlugin"

type PluginYml struct {
	MinecraftVersion minecraft.Version
	PackageJSON      *project.PackageJSON
}

func (y *PluginYml) String() string {
	var lines []string

	// General plugin details
	lines = append(lines,
		fmt.Sprintf("name: %s", y.PackageJSON.Name),
		fmt.Sprintf("api-version: %s", y.MinecraftVersion.ApiVersion()),
		fmt.Sprintf("version: %s", y.PackageJSON.Version),
		fmt.Sprintf("main: %s", JAR_MAIN_CLASS),
	)
	if len(y.PackageJSON.Author) > 0 {
		lines = append(lines, fmt.Sprintf("author: %s", y.PackageJSON.Author))
	}
	if len(y.PackageJSON.Website) > 0 {
		lines = append(lines, fmt.Sprintf("website: %s", y.PackageJSON.Website))
	}
	lines = append(lines, "")

	// Add the commands
	if len(y.PackageJSON.Commands) > 0 {
		lines = append(lines, "commands:")
		for key, attrs := range y.PackageJSON.Commands {
			lines = append(lines, indent(1)+fmt.Sprintf("%s:", key))
			if attrs != nil {
				if len(attrs.Description) > 0 {
					lines = append(lines, indent(2)+fmt.Sprintf("description: %s", attrs.Description))
				}
				if len(attrs.Aliases) > 0 {
					lines = append(lines, indent(2)+fmt.Sprintf("aliases: [%s]", strings.Join(attrs.Aliases, ", ")))
				}
				if len(attrs.Permission) > 0 {
					lines = append(lines, indent(2)+fmt.Sprintf("permission: %s", attrs.Permission))
				}
				if len(attrs.PermissionMessage) > 0 {
					lines = append(lines, indent(2)+fmt.Sprintf("permision-message: %s", attrs.PermissionMessage))
				}
				if len(attrs.Usage) > 0 {
					lines = append(lines, indent(2)+fmt.Sprintf("usage: %q", attrs.Usage))
				}
			}
		}
		lines = append(lines, "")
	}

	// Add the permissions
	if len(y.PackageJSON.Permissions) > 0 {
		lines = append(lines, "permissions:")
		for key, attrs := range y.PackageJSON.Permissions {
			lines = append(lines, indent(1)+fmt.Sprintf("%s:", key))
			if attrs != nil {
				if len(attrs.Description) > 0 {
					lines = append(lines, indent(2)+fmt.Sprintf("description: %s", attrs.Description))
				}
				if attrs.Default != nil {
					lines = append(lines, indent(2)+fmt.Sprintf("default: %t", *attrs.Default))
				}
				if attrs.Children != nil {
					lines = append(lines, indent(2)+"children:")
					for childKey, childVal := range attrs.Children {
						lines = append(lines, indent(3)+fmt.Sprintf("%s: %t", childKey, childVal))
					}
				}
			}
		}
		lines = append(lines, "")
	}

	return strings.Join(lines, "\n")
}

func indent(level int) string {
	return strings.Repeat(" ", 2*level)
}
