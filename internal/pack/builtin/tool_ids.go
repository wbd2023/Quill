package builtin

import "ciphera/tools/internal/toolchain"

const (
	ToolGo           = "go"
	ToolGoimports    = "goimports"
	ToolMisspell     = "misspell"
	ToolGolangciLint = "golangci-lint"
	ToolShfmt        = "shfmt"
	ToolShellcheck   = "shellcheck"
	ToolMarkdownlint = "markdownlint"
)

const (
	ToolVersionGoCommand  toolchain.VersionKind = "go_command"
	ToolVersionBuildInfo  toolchain.VersionKind = "build_info"
	ToolVersionShellcheck toolchain.VersionKind = "shellcheck"
	ToolVersionNodeCLI    toolchain.VersionKind = "node_cli"
)

const (
	ToolInstallNone              toolchain.InstallKind = "none"
	ToolInstallGoBinary          toolchain.InstallKind = "go_binary"
	ToolInstallNodePackage       toolchain.InstallKind = "node_package"
	ToolInstallShellcheckArchive toolchain.InstallKind = "shellcheck_archive"
)
