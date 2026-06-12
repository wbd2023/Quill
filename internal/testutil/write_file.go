package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

const (
	testDirectoryMode  = 0o700
	testExecutableMode = 0o700
	testFileMode       = 0o600
)

func WriteFile(
	test *testing.T,
	root string,
	relativePath string,
	contents string,
) (path string) {
	test.Helper()

	path = filepath.Join(root, relativePath)
	WritePath(test, path, contents)
	return path
}

func WritePath(
	test *testing.T,
	path string,
	contents string,
) {
	test.Helper()

	writePath(test, path, contents, testFileMode)
}

func WriteExecutable(
	test *testing.T,
	path string,
	contents string,
) {
	test.Helper()

	writePath(test, path, contents, testExecutableMode)
}

func ReadFile(
	test *testing.T,
	root string,
	relativePath string,
) (contents string) {
	test.Helper()

	path := filepath.Join(root, relativePath)
	data, err := os.ReadFile(path)
	if err != nil {
		test.Fatalf("read %s: %v", path, err)
	}

	return string(data)
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
