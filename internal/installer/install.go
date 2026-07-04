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

// install constants.
const (
	defaultDirectoryMode os.FileMode = 0o755
	downloadMode         os.FileMode = 0o644
	executableMode       os.FileMode = 0o755
)

/* ---------------------------------------- Installation ---------------------------------------- */

// Install downloads and installs the pinned external tools declared in the profile.
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
	switch capability.InstallKind {

	case toolchain.InstallKindNone:
		return nil

	case toolchain.InstallKindGoBinary:
		return installGoTool(layout, writer, tool, capability)

	case toolchain.InstallKindNodePackage:
		return installNodeTool(layout, writer, tool, capability)

	case toolchain.InstallKindShellcheckArchive:
		return installShellcheckTool(layout, writer, tool, capability)

	default:
		return fmt.Errorf(
			"unsupported install strategy %q for tool %s",
			capability.InstallKind,
			tool.ID,
		)
	}
}

// SupportsInstallKind reports whether kind names a known install strategy.
func SupportsInstallKind(kind toolchain.InstallKind) (supported bool) {
	switch kind {

	case toolchain.InstallKindNone,
		toolchain.InstallKindGoBinary,
		toolchain.InstallKindNodePackage,
		toolchain.InstallKindShellcheckArchive:
		return true
	}

	return false
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
