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
		InstallKind: InstallNone,
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
		VersionKind:   VersionBuildInfo,
		ModulePath:    modulePath,
		InstallKind:   InstallGoBinary,
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
		VersionKind:   VersionNodeCLI,
		InstallKind:   InstallNodePackage,
		InstallSource: installSource,
	}
}

func buildShellcheckArchive() (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:          Shellcheck,
		Name:        "shellcheck",
		Command:     "shellcheck",
		VersionKind: VersionShellcheck,
		InstallKind: InstallShellcheckArchive,
	}
}
