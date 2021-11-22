package minecraft

// SupportedVersions slice of all the supported Minecraft versions. To be supported, two things must be true:
//  1) We must have a JAR build of `bukkit-runtime` for that version
//  2) There must be a PaperMC build in that Minecraft version
var SupportedVersions = []Version{
	&paperMcVersion{"1.17.1", 386},
	&paperMcVersion{"1.16.5", 790},
}

// FindVersion finds a supported version with the given version string
func FindVersion(version string) Version {
	for _, v := range SupportedVersions {
		if v.String() == version || v.ApiVersion() == version {
			return v
		}
	}
	return nil
}

// LatestVersion gets the latest Minecraft version available
func LatestVersion() Version {
	if len(SupportedVersions) == 0 {
		return nil
	}
	return SupportedVersions[0]
}
