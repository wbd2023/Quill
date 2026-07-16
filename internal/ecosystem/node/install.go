package node

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

// Install runs npm install for the tool using an isolated npm environment derived from layout. It
// skips installation when the tool is already present at the pinned version.
func Install(
	layout runtime.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	path string,
) (err error) {
	binary := filepath.Join(BinaryDirectory(layout), tool.Command)
	installed, err := toolchain.IsInstalled(runtime.Runner{}, tool, binary)
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

	if _, err = runtime.RunCommand(command); err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

// command builds the CommandRequest for running npm install with an isolated npm environment.
// --ignore-scripts prevents arbitrary postinstall scripts from running; the remaining flags pin
// the exact version and suppress side output.
func command(
	layout runtime.Layout,
	tool toolchain.Tool,
	path string,
) (cmd runtime.CommandRequest, err error) {
	install, ok := tool.Install.(toolchain.NpmInstall)
	if !ok {
		return cmd, fmt.Errorf("tool %s is not an NPM install", tool.ID)
	}

	if install.Source == "" {
		return cmd, fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	return runtime.CommandRequest{
		Name: "npm",
		Arguments: []string{
			"install",
			"--save-exact", "--ignore-scripts", "--no-audit", "--no-fund",
			install.Source + "@" + tool.PinnedVersion,
		},
		Environment: Environment(layout, path),
		Directory:   InstallDirectory(layout),
	}, nil
}
