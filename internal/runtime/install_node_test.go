package runtime

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

func TestNpmInstallUsesLockedCiArguments(t *testing.T) {
	expected := []string{"ci", "--ignore-scripts", "--no-audit", "--no-fund"}
	if !reflect.DeepEqual(npmInstallArguments(), expected) {
		t.Fatalf("npm install arguments = %v", npmInstallArguments())
	}
}

func TestValidatePackageLockRejectsProfileVersionMismatch(t *testing.T) {
	lockPath := filepath.Join(t.TempDir(), "package-lock.json")
	contents := `{
  "packages": {
    "node_modules/markdownlint-cli": {
      "version": "0.44.0"
    }
  }
}`
	if err := os.WriteFile(lockPath, []byte(contents), downloadMode); err != nil {
		t.Fatalf("write package lock: %v", err)
	}

	err := validatePackageLock(lockPath, contract.Tool{
		PinnedVersion: "0.45.0",
	}, toolchain.Capability{
		InstallSource: "markdownlint-cli",
	})
	if err == nil || !strings.Contains(err.Error(), "profile pins 0.45.0") {
		t.Fatalf("expected lock mismatch error, got %v", err)
	}
}

func TestValidatePackageMetadataRejectsLockNameMismatch(t *testing.T) {
	directory := t.TempDir()
	packagePath := filepath.Join(directory, "package.json")
	lockPath := filepath.Join(directory, "package-lock.json")

	if err := os.WriteFile(
		packagePath,
		[]byte(`{"name":"style-platform-node-tools"}`),
		downloadMode,
	); err != nil {
		t.Fatalf("write package: %v", err)
	}

	if err := os.WriteFile(
		lockPath,
		[]byte(`{"name":"other-style-node-tools"}`),
		downloadMode,
	); err != nil {
		t.Fatalf("write package lock: %v", err)
	}

	err := validatePackageMetadata(packagePath, lockPath)
	if err == nil || !strings.Contains(err.Error(), "does not match package.json") {
		t.Fatalf("expected lock name mismatch, got %v", err)
	}
}

func TestValidatePackageMetadataAcceptsMatchingPackageAndLockNames(t *testing.T) {
	directory := t.TempDir()
	packagePath := filepath.Join(directory, "package.json")
	lockPath := filepath.Join(directory, "package-lock.json")

	err := os.WriteFile(packagePath, []byte(`{"name":"local-style-tools"}`), downloadMode)
	if err != nil {
		t.Fatalf("write package: %v", err)
	}

	err = os.WriteFile(lockPath, []byte(`{"name":"local-style-tools"}`), downloadMode)
	if err != nil {
		t.Fatalf("write package lock: %v", err)
	}

	if err := validatePackageMetadata(packagePath, lockPath); err != nil {
		t.Fatalf("validatePackageMetadata: %v", err)
	}
}
