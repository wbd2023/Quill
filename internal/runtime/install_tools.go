package runtime

import (
	"fmt"
	"io"
	"os"

	"ciphera/tools/internal/contract"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	defaultDirectoryMode os.FileMode = 0o755
	defaultDownloadMode  os.FileMode = 0o644
)

/* ---------------------------------------- Installation ---------------------------------------- */

func Install(layout Layout, writer io.Writer, tools []contract.Tool) (err error) {
	if err = ensureLayout(layout); err != nil {
		return err
	}

	for _, tool := range tools {
		if err = installTool(layout, writer, tool); err != nil {
			return err
		}
	}

	return nil
}

func installTool(layout Layout, writer io.Writer, tool contract.Tool) (err error) {
	switch tool.InstallKind {
	case contract.ToolInstallNone:
		return nil

	case contract.ToolInstallGoBinary:
		return installGoTool(layout, writer, tool)

	case contract.ToolInstallNodePackage:
		return installNodeTool(layout, writer, tool)

	case contract.ToolInstallShellcheckArchive:
		return installShellcheckTool(layout, writer, tool)

	default:
		return fmt.Errorf(
			"unsupported install strategy %q for tool %s",
			tool.InstallKind,
			tool.ID,
		)
	}
}

/* ---------------------------------------- Layout Setup ---------------------------------------- */

func ensureLayout(layout Layout) (err error) {
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
