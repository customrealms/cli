package serve

import "fmt"

type PaperVersion struct {
	Version string
	Build   int
}

var PaperVersions = []PaperVersion{
	{"1.17.1", 386},
	{"1.16.5", 790},
}

func (v *PaperVersion) Url() string {
	return fmt.Sprintf(
		"https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar",
		v.Version,
		v.Build,
		v.Version,
		v.Build,
	)
}

func FindPaperVersion(version string) *PaperVersion {
	for _, v := range PaperVersions {
		if v.Version == version {
			return &v
		}
	}
	return nil
}
