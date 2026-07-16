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

// Install runs npm install for the tool using an isolated NPM environment derived from layout. It
// skips installation when the tool is already present at the pinned version.
func Install(
	layout runtime.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	install toolchain.NpmInstall,
	path string,
) (err error) {
	if install.Source == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	binary := filepath.Join(BinaryDirectory(layout), tool.Command)
	installed, err := toolchain.IsInstalled(tool, binary)
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

	if err = os.MkdirAll(Directory(layout), standardPermissions); err != nil {
		return err
	}

	_, err = runtime.RunCommand(runtime.CommandRequest{
		Directory:   Directory(layout),
		Environment: Environment(layout, path),
		Name:        "npm",
		Arguments:   npmArguments(install.Source, tool.PinnedVersion),
	})
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

// npmArguments builds the arguments for npm install. --ignore-scripts prevents arbitrary
// postinstall scripts from running; the rest pin the exact version and suppress side output.
func npmArguments(source string, version string) (arguments []string) {
	return []string{
		"install",
		"--save-exact",
		"--ignore-scripts",
		"--no-audit",
		"--no-fund",
		source + "@" + version,
	}
}
