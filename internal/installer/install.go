package installer

import (
	"fmt"
	"io"
	"os"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	defaultDirectoryMode os.FileMode = 0o755
	downloadMode         os.FileMode = 0o644
	executableMode       os.FileMode = 0o755
)

const (
	toolInstallNone              toolchain.InstallKind = "none"
	toolInstallGoBinary          toolchain.InstallKind = "go_binary"
	toolInstallNodePackage       toolchain.InstallKind = "node_package"
	toolInstallShellcheckArchive toolchain.InstallKind = "shellcheck_archive"
)

type installHandler func(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) error

/* ---------------------------------------- Installation ---------------------------------------- */

func Install(
	layout runtime.Layout,
	writer io.Writer,
	tools []style.Tool,
	capabilities map[string]toolchain.Capability,
) (err error) {
	if err = ensureLayout(layout); err != nil {
		return err
	}

	for _, tool := range tools {
		capability, found := capabilities[tool.ID]
		if !found {
			return fmt.Errorf("missing tool capability %q", tool.ID)
		}

		if err = installTool(layout, writer, tool, capability); err != nil {
			return err
		}
	}

	return nil
}

func installTool(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	handler, found := installHandlers()[capability.InstallKind]
	if !found {
		return fmt.Errorf(
			"unsupported install strategy %q for tool %s",
			capability.InstallKind,
			tool.ID,
		)
	}

	return handler(layout, writer, tool, capability)
}

func installHandlers() (handlers map[toolchain.InstallKind]installHandler) {
	return map[toolchain.InstallKind]installHandler{
		toolInstallNone:              skipInstall,
		toolInstallGoBinary:          installGoTool,
		toolInstallNodePackage:       installNodeTool,
		toolInstallShellcheckArchive: installShellcheckTool,
	}
}

func skipInstall(
	_ runtime.Layout,
	_ io.Writer,
	_ style.Tool,
	_ toolchain.Capability,
) (err error) {
	return nil
}

func SupportsInstallKind(kind toolchain.InstallKind) (supported bool) {
	_, supported = installHandlers()[kind]
	return supported
}

/* ---------------------------------------- Layout Setup ---------------------------------------- */

func ensureLayout(layout runtime.Layout) (err error) {
	for _, path := range []string{
		layout.GoBuildCache,
		layout.GoModCache,
		layout.GoPath,
		layout.NpmCache,
		layout.ToolBinDir,
		layout.NodeBinDir,
	} {
		if err = os.MkdirAll(path, defaultDirectoryMode); err != nil {
			return err
		}
	}

	return nil
}
