package tool

import "ciphera/tools/internal/toolchain"

// install_kinds constants.
const (
	InstallNone              toolchain.InstallKind = "none"
	InstallGoBinary          toolchain.InstallKind = "go_binary"
	InstallNodePackage       toolchain.InstallKind = "node_package"
	InstallShellcheckArchive toolchain.InstallKind = "shellcheck_archive"
)
