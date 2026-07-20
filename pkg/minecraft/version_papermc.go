package minecraft

import (
	"strings"
)

type paperMcVersion struct {
	version      string
	paperBuild   int
	serverJarUrl string
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
	return v.serverJarUrl
}
