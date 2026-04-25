package runtime

import (
	"fmt"
	"io"
	"path/filepath"

	"ciphera/tools/internal/contract"
)

func installGoTool(layout Layout, writer io.Writer, tool contract.Tool) (err error) {
	if tool.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	localPath := filepath.Join(layout.ToolBinDir, tool.Command)
	localVersion, found, err := inspectLocalToolVersion(tool, localPath)
	if err != nil {
		return err
	}

	if found && matchesPinnedVersion(localVersion, tool.PinnedVersion) {
		return nil
	}

	if _, err = fmt.Fprintf(
		writer,
		"Installing %s@%s...\n",
		tool.Name,
		tool.PinnedVersion,
	); err != nil {
		return err
	}

	_, err = RunCommand(
		layout.ToolsDir,
		goInstallEnvironment(layout),
		"go",
		"install",
		tool.InstallSource+"@"+tool.PinnedVersion,
	)
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

func goInstallEnvironment(layout Layout) (environment map[string]string) {
	return map[string]string{
		"GOBIN":      layout.ToolBinDir,
		"GOCACHE":    layout.GoBuildCache,
		"GOMODCACHE": layout.GoModCache,
		"GOPATH":     layout.GoPath,
		"PATH":       layout.SearchPath(),
	}
}
