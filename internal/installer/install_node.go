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

func installNodePackage(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	if capability.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	path := filepath.Join(layout.NodeBinaryDirectory(), capability.Command)
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
		Arguments: npmArguments(capability.InstallSource, tool.PinnedVersion),
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
