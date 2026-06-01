package fixtures

import (
	"path/filepath"
	"runtime"
	"testing"
)

func RepositoryRoot(test *testing.T) (root string) {
	test.Helper()

	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		test.Fatal("runtime.Caller failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(callerFile), "..", "..", ".."))
}

func ToolsRoot(test *testing.T) (root string) {
	test.Helper()

	return filepath.Join(RepositoryRoot(test), "tools")
}
