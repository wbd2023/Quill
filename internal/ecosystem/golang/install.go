package golang

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/toolchain"
	"github.com/wbd2023/quill/internal/workspace"
)

const installTimeout = 30 * time.Minute

// Install runs go install for a tool that needs installation. The installer normally checks the
// managed binary before resolving bootstrap PATH; this package repeats that early return so direct
// callers also never need a host Go executable for an already-installed tool.
func Install(
	ctx context.Context,
	layout workspace.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	executable string,
	path string,
	runner toolchain.CommandRunner,
) (err error) {
	binary := BinaryPath(layout, tool.Command)
	installed, err := toolchain.IsInstalled(ctx, runner, tool, binary)
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

	command, err := command(layout, tool, executable, path)
	if err != nil {
		return err
	}

	if _, err = process.RunCommand(ctx, command); err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

// command builds the CommandRequest for running go install with an isolated Go environment. The
// executable is the absolute bootstrap-resolved go binary; path is the bootstrap PATH passed
// through to the child so it never searches repository or state directories.
func command(
	layout workspace.Layout,
	tool toolchain.Tool,
	executable string,
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
		Name:        executable,
		Arguments:   []string{"install", install.Source + "@" + tool.PinnedVersion},
		Environment: process.EnvironmentInherit,
		Variables:   environment,
		Directory:   layout.StateDirectory,
		Timeout:     installTimeout,
	}, nil
}
