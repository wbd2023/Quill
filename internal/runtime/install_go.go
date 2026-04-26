package runtime

import (
	"fmt"
	"io"
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

func installGoTool(
	layout Layout,
	writer io.Writer,
	tool contract.Tool,
	capability toolchain.Capability,
) (err error) {
	if capability.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	localPath := filepath.Join(layout.ToolBinDir, capability.Command)
	localVersion, found, err := inspectLocalToolVersion(tool, capability, localPath)
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

	goCapability := toolchain.Capability{ID: "go", Name: "Go", Command: "go"}
	goTool := contract.Tool{ID: "go", Name: "Go",
		TimeoutSeconds: tool.TimeoutSeconds, OutputLimitBytes: tool.OutputLimitBytes}
	_, err = RunToolCommand(
		layout.ToolsDir,
		goInstallEnvironment(layout),
		goTool,
		goCapability,
		"install",
		capability.InstallSource+"@"+tool.PinnedVersion,
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
