package installer

import (
	"fmt"
	"io"
	"path/filepath"

	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func installNodeTool(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	if capability.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	localPath := filepath.Join(layout.NodeBinaryDirectory(), capability.Command)
	installed, err := hasPinnedLocalTool(tool, capability, localPath)
	if err != nil {
		return err
	}

	if installed {
		return nil
	}

	if _, err = fmt.Fprintf(
		writer,
		"Installing %s@%s via npm install...\n",
		capability.InstallSource,
		tool.PinnedVersion,
	); err != nil {
		return err
	}

	npmTool := style.Tool{ID: "npm", Name: "npm",
		TimeoutSeconds: tool.TimeoutSeconds, OutputLimitBytes: tool.OutputLimitBytes}
	npmCapability := toolchain.Capability{ID: "npm", Name: "npm", Command: "npm"}
	_, err = runtime.RunToolCommand(
		layout.NodeDirectory(),
		map[string]string{
			"PATH":             layout.SearchPath(),
			"npm_config_cache": layout.NpmCache(),
		},
		npmTool,
		npmCapability,
		npmInstallArguments(capability.InstallSource, tool.PinnedVersion)...,
	)
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

// npmInstallArguments builds the argument list for `npm install`. npm creates package.json and
// package-lock.json automatically in the working directory, so no static lockfile is needed.
func npmInstallArguments(packageSource string, version string) (arguments []string) {
	return []string{
		"install",
		"--save-exact",
		"--ignore-scripts",
		"--no-audit",
		"--no-fund",
		packageSource + "@" + version,
	}
}
