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

/* ------------------------------------------ Constants ----------------------------------------- */

const nodeToolchainDir = "toolchain/npm"

/* ---------------------------------------- Node Install ---------------------------------------- */

func installNodeTool(
	layout runtime.Layout,
	writer io.Writer,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	if capability.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	localPath := filepath.Join(layout.NodeBinDir, capability.Command)
	installed, err := hasPinnedLocalTool(tool, capability, localPath)
	if err != nil {
		return err
	}

	if installed {
		return nil
	}

	if err = prepareLockedNodeInstall(layout, tool, capability); err != nil {
		return err
	}

	if _, err = fmt.Fprintf(
		writer,
		"Installing %s@%s via npm ci...\n",
		capability.InstallSource,
		tool.PinnedVersion,
	); err != nil {
		return err
	}

	npmTool := style.Tool{ID: "npm", Name: "npm",
		TimeoutSeconds: tool.TimeoutSeconds, OutputLimitBytes: tool.OutputLimitBytes}
	npmCapability := toolchain.Capability{ID: "npm", Name: "npm", Command: "npm"}
	_, err = runtime.RunToolCommand(
		layout.NodeDir,
		map[string]string{
			"PATH":             layout.SearchPath(),
			"npm_config_cache": layout.NpmCache,
		},
		npmTool,
		npmCapability,
		npmInstallArguments()...,
	)
	if err != nil {
		return fmt.Errorf("install %s: %w", tool.Name, err)
	}

	return nil
}

/* ------------------------------------------ Lockfiles ----------------------------------------- */

func prepareLockedNodeInstall(
	layout runtime.Layout,
	tool style.Tool,
	capability toolchain.Capability,
) (err error) {
	sourceDir := filepath.Join(layout.ToolsDir, nodeToolchainDir)
	packagePath := filepath.Join(sourceDir, "package.json")
	lockPath := filepath.Join(sourceDir, "package-lock.json")
	if err = validatePackageMetadata(packagePath, lockPath); err != nil {
		return err
	}

	if err = validatePackageLock(lockPath, tool, capability); err != nil {
		return err
	}

	if err = os.MkdirAll(layout.NodeDir, defaultDirectoryMode); err != nil {
		return err
	}

	for _, name := range []string{"package.json", "package-lock.json"} {
		if err = copyFile(
			filepath.Join(sourceDir, name),
			filepath.Join(layout.NodeDir, name),
			downloadMode,
		); err != nil {
			return err
		}
	}

	return nil
}

func npmInstallArguments() (arguments []string) {
	return []string{"ci", "--ignore-scripts", "--no-audit", "--no-fund"}
}
