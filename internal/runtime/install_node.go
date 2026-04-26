package runtime

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const nodeToolchainDir = "toolchain/npm"

/* ---------------------------------------- Node Install ---------------------------------------- */

func installNodeTool(
	layout Layout,
	writer io.Writer,
	tool contract.Tool,
	capability toolchain.Capability,
) (err error) {
	if capability.InstallSource == "" {
		return fmt.Errorf("tool %s does not define an install source", tool.ID)
	}

	localPath := filepath.Join(layout.NodeBinDir, capability.Command)
	localVersion, found, err := inspectLocalToolVersion(tool, capability, localPath)
	if err != nil {
		return err
	}

	if found && matchesPinnedVersion(localVersion, tool.PinnedVersion) {
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

	npmTool := contract.Tool{ID: "npm", Name: "npm",
		TimeoutSeconds: tool.TimeoutSeconds, OutputLimitBytes: tool.OutputLimitBytes}
	npmCapability := toolchain.Capability{ID: "npm", Name: "npm", Command: "npm"}
	_, err = RunToolCommand(
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

func npmInstallArguments() (arguments []string) {
	return []string{"ci", "--ignore-scripts", "--no-audit", "--no-fund"}
}

/* ------------------------------------------ Lockfiles ----------------------------------------- */

func prepareLockedNodeInstall(
	layout Layout,
	tool contract.Tool,
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

func validatePackageMetadata(packagePath string, lockPath string) (err error) {
	packageName, err := readPackageName(packagePath)
	if err != nil {
		return err
	}

	lockName, err := readPackageName(lockPath)
	if err != nil {
		return err
	}

	if lockName != packageName {
		return fmt.Errorf(
			"package-lock name %q does not match package.json name %q",
			lockName,
			packageName,
		)
	}

	return nil
}

func readPackageName(path string) (name string, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var document struct {
		Name string `json:"name"`
	}
	if err = json.Unmarshal(contents, &document); err != nil {
		return "", err
	}

	if document.Name == "" {
		return "", fmt.Errorf("%s does not define a package name", filepath.Base(path))
	}

	return document.Name, nil
}

func validatePackageLock(
	path string,
	tool contract.Tool,
	capability toolchain.Capability,
) (err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var document struct {
		Packages map[string]struct {
			Version string `json:"version"`
		} `json:"packages"`
	}
	if err = json.Unmarshal(contents, &document); err != nil {
		return err
	}

	packageEntry, found := document.Packages["node_modules/"+capability.InstallSource]
	if !found {
		return fmt.Errorf("package lock does not contain %s", capability.InstallSource)
	}

	if packageEntry.Version != tool.PinnedVersion {
		return fmt.Errorf(
			"package lock pins %s@%s, profile pins %s",
			capability.InstallSource,
			packageEntry.Version,
			tool.PinnedVersion,
		)
	}

	return nil
}

/* ------------------------------------------- Copying ------------------------------------------ */

func copyFile(source string, destination string, mode os.FileMode) (err error) {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := sourceFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	destinationFile, err := os.OpenFile(destination, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer func() {
		closeErr := destinationFile.Close()
		if err == nil {
			err = closeErr
		}
	}()

	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}

	return destinationFile.Chmod(mode)
}
