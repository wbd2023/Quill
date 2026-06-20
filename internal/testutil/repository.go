package testutil

import (
	"path/filepath"
	"runtime"
	"testing"
)

// RepositoryRoot repository root.
func RepositoryRoot(test *testing.T) (root string) {
	test.Helper()

	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		test.Fatal("runtime.Caller failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(callerFile), "..", "..", ".."))
}

// ToolsRoot tools root.
func ToolsRoot(test *testing.T) (root string) {
	test.Helper()

	return filepath.Join(RepositoryRoot(test), "tools")
}
