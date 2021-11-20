package papermc

import "fmt"

type Version struct {
	Version string
	Build   int
}

var SupportedVersions = []Version{
	{"1.17.1", 386},
	{"1.16.5", 790},
}

func (v *Version) Url() string {
	return fmt.Sprintf(
		"https://papermc.io/api/v2/projects/paper/versions/%s/builds/%d/downloads/paper-%s-%d.jar",
		v.Version,
		v.Build,
		v.Version,
		v.Build,
	)
}

func FindVersion(version string) *Version {
	for _, v := range SupportedVersions {
		if v.Version == version {
			return &v
		}
	}
	return nil
}
