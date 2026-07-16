package golang

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

// standardPermissions is the filesystem mode for created directories.
const standardPermissions os.FileMode = 0o755

// Install runs go install for the tool using an isolated Go environment derived from layout. It
// skips installation when the tool is already present at the pinned version.
func Install(
	layout runtime.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	install toolchain.GoInstall,
	searchPath string,
) (err error) {
	if install.Source == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	path := filepath.Join(layout.ToolBinaryDirectory(), tool.Command)
	installed, err := toolchain.IsInstalled(tool, path)
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

	environment := Environment(layout, searchPath)
	environment["GOBIN"] = layout.ToolBinaryDirectory()

	_, err = runtime.RunCommand(runtime.CommandRequest{
		Directory:   layout.StateDirectory(),
		Environment: environment,
		Name:        "go",
		Arguments:   []string{"install", install.Source + "@" + tool.PinnedVersion},
	})
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}
