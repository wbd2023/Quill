package profile

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func writeFile(t *testing.T, root string, relativePath string, contents string) {
	t.Helper()

	absolutePath := filepath.Join(root, relativePath)
	if err := os.MkdirAll(filepath.Dir(absolutePath), 0o700); err != nil {
		t.Fatalf("mkdir %s: %v", relativePath, err)
	}

	if err := os.WriteFile(absolutePath, []byte(contents), 0o600); err != nil {
		t.Fatalf("write %s: %v", relativePath, err)
	}
}

func projectRoot(t *testing.T) (root string) {
	t.Helper()

	_, callerFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}

	return filepath.Clean(filepath.Join(filepath.Dir(callerFile), "..", "..", ".."))
}
