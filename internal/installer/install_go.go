package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func installGoBinary(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	if capability.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	path := filepath.Join(layout.ToolBinaryDirectory(), capability.Command)
	installed, err := hasPinnedLocalTool(tool, capability, path)
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

	if err = os.MkdirAll(layout.StateDirectory(), standardPermissions); err != nil {
		return err
	}

	_, err = runtime.RunCommand(runtime.CommandRequest{
		Directory:   layout.StateDirectory(),
		Environment: goEnvironment(layout),
		Name:        "go",
		Arguments:   []string{"install", capability.InstallSource + "@" + tool.PinnedVersion},
	})
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

func goEnvironment(layout runtime.Layout) (environment map[string]string) {
	return map[string]string{
		"GOBIN":      layout.ToolBinaryDirectory(),
		"GOCACHE":    layout.GoBuildCache(),
		"GOMODCACHE": layout.GoModuleCache(),
		"GOPATH":     layout.GoPath(),
		"PATH":       layout.SearchPath(),
	}
}
