package rulepack

import "ciphera/tools/internal/contract"

/* ---------------------------------------- Tool Builders --------------------------------------- */

func builtinTool(
	id string,
	name string,
	command string,
	pinnedVersion string,
	versionKind contract.ToolVersionKind,
) (tool contract.Tool) {
	return contract.Tool{
		ID:            id,
		Name:          name,
		Command:       command,
		PinnedVersion: pinnedVersion,
		VersionKind:   versionKind,
		InstallKind:   contract.ToolInstallNone,
	}
}

func goBinaryTool(
	id string,
	name string,
	command string,
	pinnedVersion string,
	modulePath string,
	installSource string,
) (tool contract.Tool) {
	return contract.Tool{
		ID:            id,
		Name:          name,
		Command:       command,
		PinnedVersion: pinnedVersion,
		VersionKind:   contract.ToolVersionBuildInfo,
		ModulePath:    modulePath,
		InstallKind:   contract.ToolInstallGoBinary,
		InstallSource: installSource,
	}
}

func nodePackageTool(
	id string,
	name string,
	command string,
	pinnedVersion string,
	installSource string,
) (tool contract.Tool) {
	return contract.Tool{
		ID:            id,
		Name:          name,
		Command:       command,
		PinnedVersion: pinnedVersion,
		VersionKind:   contract.ToolVersionNodeCLI,
		InstallKind:   contract.ToolInstallNodePackage,
		InstallSource: installSource,
	}
}

func shellcheckArchiveTool() (tool contract.Tool) {
	return contract.Tool{
		ID:            contract.ToolShellcheck,
		Name:          "shellcheck",
		Command:       "shellcheck",
		PinnedVersion: "0.10.0",
		VersionKind:   contract.ToolVersionShellcheck,
		InstallKind:   contract.ToolInstallShellcheckArchive,
	}
}
