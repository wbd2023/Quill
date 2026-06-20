package tool

import "ciphera/tools/internal/toolchain"

// version_kinds constants.
const (
	VersionGoCommand  toolchain.VersionKind = "go_command"
	VersionBuildInfo  toolchain.VersionKind = "build_info"
	VersionShellcheck toolchain.VersionKind = "shellcheck"
	VersionNodeCLI    toolchain.VersionKind = "node_cli"
)
