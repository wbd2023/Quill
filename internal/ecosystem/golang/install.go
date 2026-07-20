package golang

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/wbd2023/Quill/internal/process"
	"github.com/wbd2023/Quill/internal/toolchain"
	"github.com/wbd2023/Quill/internal/workspace"
)

// standardPermissions is the filesystem mode for created directories.
const standardPermissions os.FileMode = 0o755

// Install runs go install for the tool using an isolated Go environment derived from layout. It
// skips installation when the tool is already present at the pinned version.
func Install(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	path string,
) (err error) {
	binary := filepath.Join(layout.BinaryDirectory(), tool.Command)
	installed, err := toolchain.IsInstalled(ctx, process.Runner{}, tool, binary)
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

	command, err := command(layout, tool, path)
	if err != nil {
		return err
	}

	if err = os.MkdirAll(command.Directory, standardPermissions); err != nil {
		return err
	}

	if _, err = process.RunCommand(ctx, command); err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

// command builds the CommandRequest for running go install with an isolated Go environment.
func command(
	layout workspace.Layout,
	tool toolchain.Tool,
	path string,
) (command process.CommandRequest, err error) {
	install, ok := tool.Install.(toolchain.GoInstall)
	if !ok {
		return command, fmt.Errorf("tool %s is not a Go install", tool.ID)
	}

	if install.Source == "" {
		return command, fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	environment := Environment(layout, path)
	environment["GOBIN"] = layout.BinaryDirectory()

	return process.CommandRequest{
		Name:        "go",
		Arguments:   []string{"install", install.Source + "@" + tool.PinnedVersion},
		Environment: environment,
		Directory:   layout.StateDirectory,
	}, nil
}
