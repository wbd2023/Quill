package style_test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

const localBinRelativePath = ".local/bin"
const goBinRelativePath = "bin"
const fakeGoPathDirectoryName = "gopath"
const styleRegistryTableHeader = "# tier|level|rule|name|scope|runner|target"

var scriptLibraryFiles = []string{
	"common.sh",
	"registry-constants.sh",
	"registry.sh",
	"runner.sh",
	"registry.table",
}

type scriptHarness struct {
	projectRoot       string
	styleDirectory    string
	fakeBinDirectory  string
	fakeHomeDirectory string
	fakeGoPath        string
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func newScriptHarness(t *testing.T, scriptNames ...string) (harness scriptHarness) {
	t.Helper()

	harness.projectRoot = t.TempDir()
	harness.styleDirectory = filepath.Join(harness.projectRoot, "tools", "style")
	harness.fakeBinDirectory = filepath.Join(harness.projectRoot, "fake-bin")
	harness.fakeHomeDirectory = filepath.Join(harness.projectRoot, "home")
	harness.fakeGoPath = filepath.Join(harness.projectRoot, fakeGoPathDirectoryName)

	internalDirectory := filepath.Join(harness.styleDirectory, "internal")
	entrypointsDirectory := filepath.Join(harness.styleDirectory, "entrypoints")
	if err := os.MkdirAll(entrypointsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir style entrypoints: %v", err)
	}
	if err := os.MkdirAll(internalDirectory, 0o700); err != nil {
		t.Fatalf("mkdir style internal: %v", err)
	}

	for _, scriptName := range scriptNames {
		sourcePath := filepath.Join(currentStyleDirectory(), scriptName)
		targetPath := filepath.Join(harness.styleDirectory, scriptName)
		copyFile(t, sourcePath, targetPath)
	}

	for _, fileName := range scriptLibraryFiles {
		sourcePath := filepath.Join(currentStyleDirectory(), "internal", fileName)
		targetPath := filepath.Join(internalDirectory, fileName)
		copyFile(t, sourcePath, targetPath)
	}

	testChecksDirectory := filepath.Join(harness.styleDirectory, "checks", "test")
	if err := os.MkdirAll(testChecksDirectory, 0o700); err != nil {
		t.Fatalf("mkdir style test checks: %v", err)
	}

	if err := os.MkdirAll(harness.fakeBinDirectory, 0o700); err != nil {
		t.Fatalf("mkdir fake bin: %v", err)
	}

	if err := os.MkdirAll(
		filepath.Join(harness.fakeHomeDirectory, localBinRelativePath),
		0o700,
	); err != nil {
		t.Fatalf("mkdir fake home local bin: %v", err)
	}

	if err := os.MkdirAll(filepath.Join(harness.fakeGoPath, goBinRelativePath), 0o700); err != nil {
		t.Fatalf("mkdir fake gopath bin: %v", err)
	}

	return harness
}

func (harness scriptHarness) scriptPath(scriptName string) (path string) {
	return filepath.Join(harness.styleDirectory, scriptName)
}

func (harness scriptHarness) writeProjectFile(
	t *testing.T,
	relativePath string,
	contents string,
	mode os.FileMode,
) (path string) {
	t.Helper()

	path = filepath.Join(harness.projectRoot, relativePath)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("mkdir parent for %s: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(contents), mode); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}

	return path
}

func (harness scriptHarness) writeScript(
	t *testing.T,
	scriptName string,
	contents string,
) (path string) {
	t.Helper()

	path = filepath.Join(harness.styleDirectory, scriptName)
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatalf("mkdir script parent for %s: %v", path, err)
	}

	if err := os.WriteFile(path, []byte(contents), 0o700); err != nil {
		t.Fatalf("write script %s: %v", path, err)
	}

	return path
}

func (harness scriptHarness) writeFakeCommand(
	t *testing.T,
	commandName string,
	contents string,
) (path string) {
	t.Helper()

	path = filepath.Join(harness.fakeBinDirectory, commandName)
	if err := os.WriteFile(path, []byte(contents), 0o700); err != nil {
		t.Fatalf("write fake command %s: %v", path, err)
	}

	return path
}

func (harness scriptHarness) writeProxyCommand(t *testing.T, commandName string) {
	t.Helper()

	targetPath, err := exec.LookPath(commandName)
	if err != nil {
		t.Skipf("%s is unavailable in test environment", commandName)
	}

	proxyContents := fmt.Sprintf(
		"#!/bin/bash\nset -euo pipefail\nexec %q \"$@\"\n",
		targetPath,
	)
	harness.writeFakeCommand(t, commandName, proxyContents)
}

func (harness scriptHarness) env(pathEntries []string, extra ...string) (environment []string) {
	environment = append(os.Environ(), extra...)
	environment = append(
		environment,
		"HOME="+harness.fakeHomeDirectory,
		"FAKE_GOPATH="+harness.fakeGoPath,
		"PATH="+strings.Join(pathEntries, string(os.PathListSeparator)),
	)

	return environment
}

func currentStyleDirectory() (directory string) {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), ".."))
}

func registryTableLines(rows ...string) (lines []string) {
	lines = append([]string{styleRegistryTableHeader}, rows...)
	return lines
}

func runBashScript(scriptPath string, args ...string) (output string, err error) {
	commandArguments := append([]string{scriptPath}, args...)
	command := exec.Command("bash", commandArguments...)

	rawOutput, runErr := command.CombinedOutput()
	return string(rawOutput), runErr
}

func runBashScriptWithEnv(
	scriptPath string,
	environment []string,
	args ...string,
) (output string, err error) {
	commandArguments := append([]string{scriptPath}, args...)
	command := exec.Command("bash", commandArguments...)
	command.Env = environment

	rawOutput, runErr := command.CombinedOutput()
	return string(rawOutput), runErr
}

func copyFile(t *testing.T, sourcePath string, targetPath string) {
	t.Helper()

	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("read %s: %v", sourcePath, err)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0o700); err != nil {
		t.Fatalf("mkdir parent for %s: %v", targetPath, err)
	}

	if err := os.WriteFile(targetPath, contents, 0o600); err != nil {
		t.Fatalf("write %s: %v", targetPath, err)
	}
}
