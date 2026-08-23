package node

import (
	"context"
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/toolchain"
	"github.com/wbd2023/quill/internal/workspace"
)

// Install runs npm install for a tool that needs installation. The installer normally checks the
// managed binary before resolving bootstrap PATH; this package repeats that early return so direct
// callers also never need a host npm executable for an already-installed tool.
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

// command builds the CommandRequest for running npm install with an isolated npm environment.
// --ignore-scripts prevents arbitrary postinstall scripts; --no-save and --package-lock=false
// prevent npm from rewriting repository-state manifests. The executable is the trusted,
// bootstrap-resolved npm binary, and path never searches repository or state directories.
func command(
	layout workspace.Layout,
	tool toolchain.Tool,
	executable string,
	path string,
) (cmd process.CommandRequest, err error) {
	install, ok := tool.Install.(toolchain.NpmInstall)
	if !ok {
		return cmd, fmt.Errorf("tool %s is not an NPM install", tool.ID)
	}

	if install.Source == "" {
		return cmd, fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	return process.CommandRequest{
		Name: executable,
		Arguments: []string{
			"install",
			"--ignore-scripts", "--no-save", "--package-lock=false", "--no-audit", "--no-fund",
			install.Source + "@" + tool.PinnedVersion,
		},
		Environment: process.EnvironmentInherit,
		Variables:   Environment(layout, path),
		Directory:   InstallDirectory(layout),
	}, nil
}
