package runtime

import (
	"fmt"
	"io"
	"path/filepath"

	"ciphera/tools/internal/contract"
)

func installNodeTool(layout Layout, writer io.Writer, tool contract.Tool) (err error) {
	if tool.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	localPath := filepath.Join(layout.NodeBinDir, tool.Command)
	localVersion, found, err := inspectLocalToolVersion(tool, localPath)
	if err != nil {
		return err
	}

	if found && matchesPinnedVersion(localVersion, tool.PinnedVersion) {
		return nil
	}

	if _, err = fmt.Fprintf(
		writer,
		"Installing %s@%s via npm prefix...\n",
		tool.InstallSource,
		tool.PinnedVersion,
	); err != nil {
		return err
	}

	_, err = RunCommand(
		layout.ToolsDir,
		map[string]string{
			"PATH":             layout.SearchPath(),
			"npm_config_cache": layout.NpmCache,
		},
		"npm",
		"install",
		"--prefix",
		layout.NodeDir,
		tool.InstallSource+"@"+tool.PinnedVersion,
	)
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}
