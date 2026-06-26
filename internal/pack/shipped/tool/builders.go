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

func buildShellcheckArchive() (capability toolchain.Capability) {
	return toolchain.Capability{
		ID:          Shellcheck,
		Name:        "shellcheck",
		Command:     "shellcheck",
		VersionKind: toolchain.VersionKindShellcheck,
		InstallKind: toolchain.InstallKindShellcheckArchive,
	}
}
