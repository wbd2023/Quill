package golang

import (
	"testing"

	"ciphera/tools/internal/runtime"
)

func TestBuildCacheDirectoryDerivesFromLayout(t *testing.T) {
	layout := runtime.NewLayout("/repo")
	expected := "/repo/.cache/quill/cache/go-build"
	if actual := BuildCacheDirectory(layout); actual != expected {
		t.Fatalf("BuildCacheDirectory = %q, want %q", actual, expected)
	}
}
