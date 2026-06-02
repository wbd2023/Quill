package builtin

import (
	"ciphera/tools/internal/pack/builtin/bash"
	"ciphera/tools/internal/pack/builtin/golang"
	"ciphera/tools/internal/pack/builtin/markdown"
	"ciphera/tools/internal/pack/builtin/text"
	"ciphera/tools/internal/toolchain"
)

const (
	ToolGo           = golang.ToolGo
	ToolGoimports    = golang.ToolGoimports
	ToolMisspell     = text.ToolMisspell
	ToolGolangciLint = golang.ToolGolangciLint
	ToolShfmt        = bash.ToolShfmt
	ToolShellcheck   = bash.ToolShellcheck
	ToolMarkdownlint = markdown.ToolMarkdownlint
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
