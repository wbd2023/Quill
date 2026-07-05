package installer

import (
	"fmt"
	"io"
	"path/filepath"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func installGoTool(
	layout runtime.Layout,
	toolsDirectory string,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	if capability.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	localPath := filepath.Join(layout.ToolBinaryDirectory(), capability.Command)
	installed, err := hasPinnedLocalTool(tool, capability, localPath)
	if err != nil {
		return err
	}

	if installed {
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
	goTool := style.Tool{ID: "go", Name: "Go",
		TimeoutSeconds: tool.TimeoutSeconds, OutputLimitBytes: tool.OutputLimitBytes}
	_, err = runtime.RunToolCommand(
		toolsDirectory,
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

func goInstallEnvironment(layout runtime.Layout) (environment map[string]string) {
	return map[string]string{
		"GOBIN":      layout.ToolBinaryDirectory(),
		"GOCACHE":    layout.GoBuildCache(),
		"GOMODCACHE": layout.GoModuleCache(),
		"GOPATH":     layout.GoPath(),
		"PATH":       layout.SearchPath(),
	}
}
