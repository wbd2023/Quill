package tool

import "ciphera/tools/internal/toolchain"

func buildBuiltin(
	id string,
	name string,
	command string,
	versionKind toolchain.VersionKind,
) (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:          id,
		Name:        name,
		Command:     command,
		VersionKind: versionKind,
		InstallKind: toolchain.InstallKindNone,
	}
}

func buildGoBinary(
	id string,
	name string,
	command string,
	modulePath string,
	installSource string,
) (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:            id,
		Name:          name,
		Command:       command,
		VersionKind:   toolchain.VersionKindBuildInfo,
		ModulePath:    modulePath,
		InstallKind:   toolchain.InstallKindGoBinary,
		InstallSource: installSource,
	}
}

func buildNodePackage(
	id string,
	name string,
	command string,
	installSource string,
) (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:            id,
		Name:          name,
		Command:       command,
		VersionKind:   toolchain.VersionKindNodeCLI,
		InstallKind:   toolchain.InstallKindNodePackage,
		InstallSource: installSource,
	}
}

func shellcheckArchiveURL(version string, platform string) (url string) {
	return "https://github.com/koalaman/shellcheck/releases/download/v" +
		version + "/shellcheck-v" + version + "." + platform + ".tar.xz"
}

func shellcheckBinaryPath(version string) (path string) {
	return "shellcheck-v" + version + "/shellcheck"
}

func buildShellcheckArchive() (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:          Shellcheck,
		Name:        "shellcheck",
		Command:     "shellcheck",
		VersionKind: toolchain.VersionKindShellcheck,
		InstallKind: toolchain.InstallKindArchive,
		Archive: &toolchain.ArchiveSpec{
			URL:        shellcheckArchiveURL,
			Format:     toolchain.ArchiveFormatXz,
			BinaryPath: shellcheckBinaryPath,
			Platforms: map[string]string{
				"darwin/amd64": "darwin.x86_64",
				"darwin/arm64": "darwin.aarch64",
				"linux/amd64":  "linux.x86_64",
				"linux/arm64":  "linux.aarch64",
			},
		},
	}
}
