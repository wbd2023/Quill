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

/* ---------------------------------------- File Writing ---------------------------------------- */

// WriteFile creates a file at root/path with owner-only permissions (0o600), creating parent
// directories as needed. Returns the joined absolute path. Use WriteFileAt when you already have
// the absolute path.
func WriteFile(
	test *testing.T,
	root string,
	path string,
	contents string,
) (absolutePath string) {
	test.Helper()

	absolutePath = filepath.Join(root, path)
	WriteFileAt(test, absolutePath, contents)
	return absolutePath
}

// WriteFileAt creates a file at the given absolute path with owner-only permissions (0o600),
// creating parent directories as needed.
func WriteFileAt(
	test *testing.T,
	path string,
	contents string,
) {
	test.Helper()

	writePath(test, path, contents, testFileMode)
}

// WriteExecutable creates an executable file at the given absolute path (0o755), creating parent
// directories as needed. Used for test fixtures that simulate installed binaries.
func WriteExecutable(
	test *testing.T,
	path string,
	contents string,
) {
	test.Helper()

	writePath(test, path, contents, testExecutableMode)
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

/* ---------------------------------------- File Reading ---------------------------------------- */

// ReadFile reads the file at root/path and returns its contents.
func ReadFile(
	test *testing.T,
	root string,
	path string,
) (contents string) {
	test.Helper()

	absolutePath := filepath.Join(root, path)
	data, err := os.ReadFile(absolutePath)
	if err != nil {
		test.Fatalf("read %s: %v", absolutePath, err)
	}

	return string(data)
}
