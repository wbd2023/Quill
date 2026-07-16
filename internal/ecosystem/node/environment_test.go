package node

import (
	"testing"

	"ciphera/tools/internal/runtime"
)

func TestInstallDirectoryDerivesFromLayout(t *testing.T) {
	layout := runtime.NewLayout("/repo")
	expected := "/repo/.cache/quill/npm"
	if actual := InstallDirectory(layout); actual != expected {
		t.Fatalf("InstallDirectory = %q, want %q", actual, expected)
	}
}

func TestCacheDirectoryDerivesFromLayout(t *testing.T) {
	layout := runtime.NewLayout("/repo")
	expected := "/repo/.cache/quill/cache/npm"
	if actual := CacheDirectory(layout); actual != expected {
		t.Fatalf("CacheDirectory = %q, want %q", actual, expected)
	}
}
