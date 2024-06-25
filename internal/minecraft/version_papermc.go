package minecraft

import (
	"fmt"
	"strings"
)

type paperMcVersion struct {
	version    string
	paperBuild int
}

func (v *paperMcVersion) String() string {
	return v.version
}

func (v *paperMcVersion) ApiVersion() string {
	parts := strings.Split(v.version, ".")
	return strings.Join(parts[:2], ".")
}

func (v *paperMcVersion) ServerJarType() string {
	return "paper"
}

func (v *paperMcVersion) ServerJarUrl() string {
	return fmt.Sprintf(
		"https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar",
		v.version,
		v.paperBuild,
		v.version,
		v.paperBuild,
	)
}
