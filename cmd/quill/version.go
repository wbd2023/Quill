package main

import "runtime/debug"

// developmentVersion is Go's version for an unversioned main module.
const developmentVersion = "(devel)"

func version() (value string) {
	return resolveVersion(debug.ReadBuildInfo())
}

func resolveVersion(info *debug.BuildInfo, ok bool) (version string) {
	if !ok || info.Main.Version == "" {
		return developmentVersion
	}

	return info.Main.Version
}
