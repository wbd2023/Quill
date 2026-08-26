package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/workspace"
)

/* ---------------------------------- State Boundary Rejections --------------------------------- */

// TestPrepareInstallDirectoryRejectsSymlinkedCacheComponent is the QUILL-TRUST-001 regression: a
// cache or state component that is a symlink is rejected before Go or npm writes follow it outside
// the repository.
func TestPrepareInstallDirectoryRejectsSymlinkedCacheComponent(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	sentinel := filepath.Join(outside, "untouched")
	if err := os.WriteFile(sentinel, []byte("protected"), 0o600); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}

	// The repository state tree is replaced by a symlink into an outside directory.
	if err := os.Symlink(outside, filepath.Join(root, ".cache")); err != nil {
		t.Fatalf("create state symlink: %v", err)
	}

	layout := workspace.NewLayout(root)
	if err := prepareGoInstall(layout); err == nil {
		t.Fatal("prepareGoInstall = nil, want error for symlinked state component")
	}

	// No write must have escaped through the symlink before the rejection.
	if content, err := os.ReadFile(sentinel); err != nil ||
		string(content) != "protected" {
		t.Fatalf("outside sentinel changed: content=%q err=%v", content, err)
	}
}

// TestPrepareInstallDirectoryRejectsCustomStateOutsideRepository closes the custom-state path of
// QUILL-TRUST-001: a state directory configured outside the repository root is rejected before
// any third-party write.
func TestPrepareInstallDirectoryRejectsCustomStateOutsideRepository(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	layout := workspace.NewLayout(root)
	layout.StateDirectory = filepath.Join(outside, "quill-state")

	if err := prepareNpmInstall(layout); err == nil {
		t.Fatal("prepareNpmInstall = nil, want error for state outside repository")
	}

	// Nothing should have been created in the outside state tree.
	if _, err := os.Stat(layout.StateDirectory); !os.IsNotExist(err) {
		t.Fatalf("outside state directory exists after rejection: %v", err)
	}
}

// TestPrepareNpmInstallRejectsManifestSymlink proves npm cannot follow a repository-state
// manifest link before it starts resolving a package installation.
func TestPrepareNpmInstallRejectsManifestSymlink(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "package.json")
	if err := os.WriteFile(outside, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write outside manifest: %v", err)
	}

	layout := workspace.NewLayout(root)
	manifest := filepath.Join(layout.StateDirectory, "npm", "package.json")
	if err := os.MkdirAll(filepath.Dir(manifest), 0o755); err != nil {
		t.Fatalf("mkdir npm state: %v", err)
	}
	if err := os.Symlink(outside, manifest); err != nil {
		t.Fatalf("create manifest symlink: %v", err)
	}

	err := prepareNpmInstall(layout)
	if err == nil || !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("prepareNpmInstall error = %v, want manifest rejection", err)
	}
}

// TestPrepareInstallDirectoryRejectsSymlinkedGoCache isolates the Go-specific cache paths
// (GOCACHE, GOMODCACHE, GOPATH) beneath the shared cache directory.
func TestPrepareInstallDirectoryRejectsSymlinkedGoCache(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()

	if err := os.MkdirAll(filepath.Join(root, ".cache", "quill"), 0o755); err != nil {
		t.Fatalf("mkdir state: %v", err)
	}
	// The shared cache directory itself is a symlink outward.
	if err := os.Symlink(outside, filepath.Join(root, ".cache", "quill", "cache")); err != nil {
		t.Fatalf("create cache symlink: %v", err)
	}

	layout := workspace.NewLayout(root)
	err := prepareGoInstall(layout)
	if err == nil {
		t.Fatal("prepareGoInstall = nil, want error for symlinked cache component")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("prepareGoInstall error = %v, want a symlink rejection", err)
	}
}

/* ---------------------------------- Preparation Success Path ---------------------------------- */

// TestPrepareGoInstallCreatesRealDirectories confirms the positive path: ordinary missing
// directories are created as real directories so normal installation is unaffected.
func TestPrepareGoInstallCreatesRealDirectories(t *testing.T) {
	root := t.TempDir()
	layout := workspace.NewLayout(root)

	if err := prepareGoInstall(layout); err != nil {
		t.Fatalf("prepareGoInstall: %v", err)
	}

	for _, directory := range []string{
		layout.StateDirectory,
		layout.BinaryDirectory(),
		layout.CacheDirectory(),
	} {
		info, err := os.Lstat(directory)
		if err != nil {
			t.Fatalf("Lstat %s: %v", directory, err)
		}
		if info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("prepared directory %s is a symlink", directory)
		}
		if !info.IsDir() {
			t.Fatalf("prepared path %s is not a directory", directory)
		}
	}
}
