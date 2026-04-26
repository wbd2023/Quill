package rulepack

import "ciphera/tools/internal/toolchain"

func builtinTool(
	id string,
	name string,
	command string,
	versionKind toolchain.VersionKind,
) (tool toolchain.Capability) {
	return toolchain.Capability{
		ID:          id,
		Name:        name,
		Command:     command,
		VersionKind: versionKind,
		InstallKind: ToolInstallNone,
	}
}

func goBinaryTool(
	id string,
	name string,
	command string,
	modulePath string,
	installSource string,
) (tool toolchain.Capability) {
	return toolchain.Capability{
		ID:            id,
		Name:          name,
		Command:       command,
		VersionKind:   ToolVersionBuildInfo,
		ModulePath:    modulePath,
		InstallKind:   ToolInstallGoBinary,
		InstallSource: installSource,
	}
}

func nodePackageTool(
	id string,
	name string,
	command string,
	installSource string,
) (tool toolchain.Capability) {
	return toolchain.Capability{
		ID:            id,
		Name:          name,
		Command:       command,
		VersionKind:   ToolVersionNodeCLI,
		InstallKind:   ToolInstallNodePackage,
		InstallSource: installSource,
	}
}

func shellcheckArchiveTool() (tool toolchain.Capability) {
	return toolchain.Capability{
		ID:          ToolShellcheck,
		Name:        "shellcheck",
		Command:     "shellcheck",
		VersionKind: ToolVersionShellcheck,
		InstallKind: ToolInstallShellcheckArchive,
	}
}
