package main

import (
	"fmt"
	"strings"
)

type PluginYml struct {
	Name       string
	ApiVersion string
	Version    string
	Main       string
}

func (y *PluginYml) String() string {
	lines := make([]string, 4)
	lines[0] = fmt.Sprintf("name: %s", y.Name)
	lines[1] = fmt.Sprintf("api-version: %s", y.ApiVersion)
	lines[2] = fmt.Sprintf("version: %s", y.Version)
	lines[3] = fmt.Sprintf("main: %s", y.Main)
	return strings.Join(lines, "\n") + "\n"
}
