package testutil

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// RepositoryRoot returns the Go module root containing the test helpers.
func RepositoryRoot(test *testing.T) (root string) {
	test.Helper()

	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		test.Fatal("runtime.Caller failed")
	}

	for root = filepath.Dir(callerFile); ; root = filepath.Dir(root) {
		_, err := os.Stat(filepath.Join(root, "go.mod"))
		switch {
		case err == nil:
			return root
		case !errors.Is(err, os.ErrNotExist):
			test.Fatalf("stat module root marker: %v", err)
		}

		parent := filepath.Dir(root)
		if parent == root {
			test.Fatal("Go module root not found")
		}
	}
}
