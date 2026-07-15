package installer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/toolchain"
)

func installNpm(
	layout runtime.Layout,
	writer io.Writer,
	tool toolchain.Tool,
	install toolchain.NpmInstall,
) (err error) {
	if install.Source == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	path := filepath.Join(layout.NodeBinaryDirectory(), tool.Command)
	installed, err := hasPinnedLocalTool(tool, path)
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

	if err = os.MkdirAll(layout.NodeDirectory(), standardPermissions); err != nil {
		return err
	}

	_, err = runtime.RunCommand(runtime.CommandRequest{
		Directory: layout.NodeDirectory(),
		Environment: map[string]string{
			"PATH":             layout.SearchPath(),
			"npm_config_cache": layout.NpmCache(),
		},
		Name:      "npm",
		Arguments: npmArguments(install.Source, tool.PinnedVersion),
	})
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

// npmArguments builds the arguments for `npm install`. --ignore-scripts prevents arbitrary
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
