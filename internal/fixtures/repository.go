package fixtures

import (
	"path/filepath"
	"runtime"
	"testing"
)

/* ------------------------------------- Current Repository ------------------------------------- */

func RepoRoot(test *testing.T) (root string) {
	test.Helper()

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		test.Fatal("runtime.Caller failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(currentFile), "..", "..", ".."))
}

func ToolsRoot(test *testing.T) (root string) {
	test.Helper()

	return filepath.Join(RepoRoot(test), "tools")
}
