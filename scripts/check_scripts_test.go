package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------------- Tests ------------------------------------------- */

func TestCheckBashLineLengthScript(t *testing.T) {
	tempProjectRoot, targetScript := setupBashScriptHarness(t, "check-bash-line-length.sh")

	shortScriptPath := filepath.Join(tempProjectRoot, "tools", "short.sh")
	shortScriptContent := "#!/bin/bash\nset -euo pipefail\necho \"short\"\n"
	if err := os.WriteFile(shortScriptPath, []byte(shortScriptContent), 0o600); err != nil {
		t.Fatalf("write short script: %v", err)
	}

	if output, err := runBashScript(targetScript, "--scope", "tools"); err != nil {
		t.Fatalf("expected short script to pass, output:\n%s", output)
	}

	tabWidthScriptPath := filepath.Join(tempProjectRoot, "tools", "tab-width.sh")
	tabWidthScriptContent := "#!/bin/bash\nset -euo pipefail\n\t" + strings.Repeat("a", 97) + "\n"
	if err := os.WriteFile(tabWidthScriptPath, []byte(tabWidthScriptContent), 0o600); err != nil {
		t.Fatalf("write tab-width script: %v", err)
	}

	output, err := runBashScript(targetScript, "--scope", "tools")
	if err == nil {
		t.Fatalf("expected tab-width script to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "[1.1] Bash line exceeds 100 columns") {
		t.Fatalf("expected tab-width line length finding, got:\n%s", output)
	}

	tabWidthAllowedContent := "#!/bin/bash\nset -euo pipefail\n\t" + strings.Repeat("a", 97)
	tabWidthAllowedContent += " # style: allow-long-line\n"
	if err := os.WriteFile(tabWidthScriptPath, []byte(tabWidthAllowedContent), 0o600); err != nil {
		t.Fatalf("rewrite tab-width script with marker: %v", err)
	}

	if output, err := runBashScript(targetScript, "--scope", "tools"); err != nil {
		t.Fatalf("expected marked tab-width line to pass, output:\n%s", output)
	}

	longScriptPath := filepath.Join(tempProjectRoot, "tools", "long.sh")
	longPayload := strings.Repeat("a", 120)
	longScriptContent := "#!/bin/bash\nset -euo pipefail\necho \"" + longPayload + "\"\n"
	if err := os.WriteFile(longScriptPath, []byte(longScriptContent), 0o600); err != nil {
		t.Fatalf("write long script: %v", err)
	}

	output, err = runBashScript(targetScript, "--scope", "tools")
	if err == nil {
		t.Fatalf("expected long script to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "[1.1] Bash line exceeds 100 columns") {
		t.Fatalf("expected line length finding, got:\n%s", output)
	}

	allowedLongScriptContent := "#!/bin/bash\nset -euo pipefail\n"
	allowedLongScriptContent += "echo \"" + longPayload + "\" # style: allow-long-line\n"
	if err := os.WriteFile(longScriptPath, []byte(allowedLongScriptContent), 0o600); err != nil {
		t.Fatalf("rewrite long script with marker: %v", err)
	}

	if output, err := runBashScript(targetScript, "--scope", "tools"); err != nil {
		t.Fatalf("expected marked long line to pass, output:\n%s", output)
	}
}

func TestCheckBashLineLengthScriptUsage(t *testing.T) {
	_, targetScript := setupBashScriptHarness(t, "check-bash-line-length.sh")

	output, err := runBashScript(targetScript, "--bad-flag")
	if err == nil {
		t.Fatalf("expected invalid arguments to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "usage:") {
		t.Fatalf("expected usage output for invalid arguments, got:\n%s", output)
	}

	output, err = runBashScript(targetScript, "--scope")
	if err == nil {
		t.Fatalf("expected missing scope value to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "usage:") {
		t.Fatalf("expected usage output for missing scope value, got:\n%s", output)
	}
}

func TestCheckGoLineLengthScriptTabWidth(t *testing.T) {
	tempProjectRoot, targetScript := setupBashScriptHarness(t, "check-go-line-length.sh")

	fakeBinDirectory := filepath.Join(tempProjectRoot, "fake-bin")
	if err := os.MkdirAll(fakeBinDirectory, 0o700); err != nil {
		t.Fatalf("mkdir fake bin directory: %v", err)
	}

	fakeLinterPath := filepath.Join(fakeBinDirectory, "golangci-lint")
	fakeLinterContent := "#!/bin/bash\nset -euo pipefail\n"
	fakeLinterContent += "if [ \"${1:-}\" = \"run\" ] && [ \"${2:-}\" = \"--help\" ]; then\n"
	fakeLinterContent += "\techo \"--enable-only\"\n\texit 0\nfi\n"
	fakeLinterContent += "if [ \"${1:-}\" = \"run\" ]; then\n\texit 0\nfi\nexit 0\n"
	if err := os.WriteFile(fakeLinterPath, []byte(fakeLinterContent), 0o700); err != nil {
		t.Fatalf("write fake golangci-lint: %v", err)
	}

	goFilePath := filepath.Join(tempProjectRoot, "tools", "tab_width.go")
	goFileContent := "package tools\n\nfunc tabWidth() {\n\t// " + strings.Repeat("a", 97) + "\n}\n"
	if err := os.WriteFile(goFilePath, []byte(goFileContent), 0o600); err != nil {
		t.Fatalf("write tab-width go file: %v", err)
	}

	goTestFilePath := filepath.Join(tempProjectRoot, "tools", "tab_width_test.go")
	goTestFileContent := "package tools\n\nfunc TestTabWidth() {\n\t// "
	goTestFileContent += strings.Repeat("a", 97) + "\n}\n"
	if err := os.WriteFile(goTestFilePath, []byte(goTestFileContent), 0o600); err != nil {
		t.Fatalf("write tab-width go test file: %v", err)
	}

	pathOverride := "PATH=" + fakeBinDirectory + ":" + os.Getenv("PATH")
	output, err := runBashScriptWithEnv(targetScript, []string{pathOverride}, "--scope", "tools")
	if err == nil {
		t.Fatalf("expected tab-width go file to fail, output:\n%s", output)
	}

	if !strings.Contains(output, "[1.1] Go line exceeds 100 columns") {
		t.Fatalf("expected tab-width finding, got:\n%s", output)
	}

	fixedGoFileContent := "package tools\n\nfunc tabWidth() {\n\t// short\n}\n"
	if err := os.WriteFile(goFilePath, []byte(fixedGoFileContent), 0o600); err != nil {
		t.Fatalf("rewrite tab-width go file: %v", err)
	}

	fixedGoTestFileContent := "package tools\n\nfunc TestTabWidth() {\n\t// short\n}\n"
	if err := os.WriteFile(goTestFilePath, []byte(fixedGoTestFileContent), 0o600); err != nil {
		t.Fatalf("rewrite tab-width go test file: %v", err)
	}

	if output, err := runBashScriptWithEnv(
		targetScript,
		[]string{pathOverride},
		"--scope",
		"tools",
	); err != nil {
		t.Fatalf("expected fixed go file to pass, output:\n%s", output)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func setupBashScriptHarness(
	t *testing.T,
	scriptName string,
) (tempProjectRoot string, targetScript string) {
	t.Helper()

	tempProjectRoot = t.TempDir()
	tempScriptsDirectory := filepath.Join(tempProjectRoot, "tools", "scripts")
	if err := os.MkdirAll(tempScriptsDirectory, 0o700); err != nil {
		t.Fatalf("mkdir scripts directory: %v", err)
	}

	tempLibraryDirectory := filepath.Join(tempScriptsDirectory, "lib")
	if err := os.MkdirAll(tempLibraryDirectory, 0o700); err != nil {
		t.Fatalf("mkdir script lib directory: %v", err)
	}

	sourceScript := filepath.Join(currentScriptsDirectory(), scriptName)
	targetScript = filepath.Join(tempScriptsDirectory, scriptName)
	copyFile(t, sourceScript, targetScript)

	sourceLibraryFile := filepath.Join(currentScriptsDirectory(), "lib", "style-common.sh")
	targetLibraryFile := filepath.Join(tempLibraryDirectory, "style-common.sh")
	copyFile(t, sourceLibraryFile, targetLibraryFile)

	return tempProjectRoot, targetScript
}

func runBashScript(scriptPath string, args ...string) (output string, err error) {
	commandArguments := append([]string{scriptPath}, args...)
	command := exec.Command("bash", commandArguments...)

	rawOutput, runErr := command.CombinedOutput()
	return string(rawOutput), runErr
}

func runBashScriptWithEnv(
	scriptPath string,
	envelope []string,
	args ...string,
) (output string, err error) {
	commandArguments := append([]string{scriptPath}, args...)
	command := exec.Command("bash", commandArguments...)
	command.Env = append(os.Environ(), envelope...)

	rawOutput, runErr := command.CombinedOutput()
	return string(rawOutput), runErr
}

func copyFile(t *testing.T, sourcePath string, targetPath string) {
	t.Helper()

	contents, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatalf("read %s: %v", sourcePath, err)
	}

	if err := os.WriteFile(targetPath, contents, 0o600); err != nil {
		t.Fatalf("write %s: %v", targetPath, err)
	}
}
