package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	testDirectoryMode  os.FileMode = 0o700
	testExecutableMode os.FileMode = 0o755
	testFileMode       os.FileMode = 0o600
)

// WriteFile creates a file at root/path with owner-only permissions (0o600), creating parent
// directories as needed. Returns the joined path.
func WriteFile(
	test *testing.T,
	root string,
	path string,
	contents string,
) (joined string) {
	test.Helper()

	joined = filepath.Join(root, path)
	writePath(test, joined, contents, testFileMode)
	return joined
}

// WriteExecutable creates an executable file at root/path with executable permissions (0o755),
// creating parent directories as needed. Used for test fixtures that simulate installed binaries.
func WriteExecutable(
	test *testing.T,
	root string,
	path string,
	contents string,
) {
	test.Helper()

	writePath(test, filepath.Join(root, path), contents, testExecutableMode)
}

func writePath(
	test *testing.T,
	path string,
	contents string,
	mode os.FileMode,
) {
	test.Helper()

	if err := os.MkdirAll(filepath.Dir(path), testDirectoryMode); err != nil {
		test.Fatalf("mkdir %s: %v", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, []byte(contents), mode); err != nil {
		test.Fatalf("write %s: %v", path, err)
	}
}

// ReadFile reads the file at root/path and returns its contents.
func ReadFile(
	test *testing.T,
	root string,
	path string,
) (contents string) {
	test.Helper()

	full := filepath.Join(root, path)
	data, err := os.ReadFile(full)
	if err != nil {
		test.Fatalf("read %s: %v", full, err)
	}

	return string(data)
}
