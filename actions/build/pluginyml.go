package build

import (
	"fmt"
	"strings"

	"github.com/customrealms/cli/project"
)

const JAR_MAIN_CLASS = "io.customrealms.MainPlugin"

type PluginYml struct {
	MinecraftVersion string
	PackageJSON      *project.PackageJSON
}

func (y *PluginYml) String() string {
	lines := make([]string, 4)
	lines[0] = fmt.Sprintf("name: %s", y.PackageJSON.Name)
	lines[1] = fmt.Sprintf("api-version: %s", mcVersionToApiVersion(y.MinecraftVersion))
	lines[2] = fmt.Sprintf("version: %s", y.PackageJSON.Version)
	lines[3] = fmt.Sprintf("main: %s", JAR_MAIN_CLASS)
	return strings.Join(lines, "\n") + "\n"
}

func mcVersionToApiVersion(mcVersion string) string {
	parts := strings.Split(mcVersion, ".")
	return strings.Join(parts[:2], ".")
}
