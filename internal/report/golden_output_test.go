package report

import (
	"os"
	"path/filepath"
	"testing"
)

func readGoldenOutput(t *testing.T, name string) (output string) {
	t.Helper()

	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden output %q: %v", name, err)
	}

	return string(data)
}
