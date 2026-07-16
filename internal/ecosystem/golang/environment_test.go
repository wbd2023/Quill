package golang

import (
	"testing"

	"ciphera/tools/internal/runtime"
)

func TestBuildCacheDerivesFromLayout(t *testing.T) {
	layout := runtime.NewLayout("/repo")
	expected := "/repo/.cache/quill/cache/go-build"
	if actual := BuildCache(layout); actual != expected {
		t.Fatalf("BuildCache = %q, want %q", actual, expected)
	}
}
