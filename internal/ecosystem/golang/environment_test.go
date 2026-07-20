package golang

import (
	"testing"

	"github.com/wbd2023/Quill/internal/workspace"
)

func TestBuildCacheDirectoryDerivesFromLayout(t *testing.T) {
	layout := workspace.NewLayout("/repo")
	expected := "/repo/.cache/quill/cache/go-build"
	if actual := BuildCacheDirectory(layout); actual != expected {
		t.Fatalf("BuildCacheDirectory = %q, want %q", actual, expected)
	}
}
