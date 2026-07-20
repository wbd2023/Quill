package main

import "runtime/debug"

const developmentVersion = "devel"

func currentVersion() (version string) {
	buildInfo, ok := debug.ReadBuildInfo()
	return versionFromBuildInfo(buildInfo, ok)
}

func versionFromBuildInfo(buildInfo *debug.BuildInfo, ok bool) (version string) {
	if !ok || buildInfo.Main.Version == "" || buildInfo.Main.Version == "(devel)" {
		return developmentVersion
	}

	return buildInfo.Main.Version
}
