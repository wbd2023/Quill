package runtime

import "testing"

func TestLayoutDerivesPathsFromRepositoryRoot(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")

	layout := NewLayout("/repo")

	if layout.StateDirectory() != "/repo/.cache/quill" {
		t.Fatalf("StateDirectory = %q, want %q", layout.StateDirectory(), "/repo/.cache/quill")
	}

	if layout.ToolBinaryDirectory() != "/repo/.cache/quill/bin" {
		t.Fatalf("ToolBinaryDirectory = %q, want %q",
			layout.ToolBinaryDirectory(), "/repo/.cache/quill/bin")
	}

	if layout.GoBuildCache() != "/repo/.cache/quill/cache/go-build" {
		t.Fatalf("GoBuildCache = %q, want %q",
			layout.GoBuildCache(), "/repo/.cache/quill/cache/go-build")
	}
}

func TestLayoutSearchPathPrependsToolAndNodeBins(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")

	layout := NewLayout("/repo")

	expected := layout.ToolBinaryDirectory() + ":" + layout.NodeBinaryDirectory() + ":/usr/bin"
	if actual := layout.SearchPath(); actual != expected {
		t.Fatalf("PATH = %q, want %q", actual, expected)
	}
}

func TestLayoutWithStateDirectoryOverridesDefault(t *testing.T) {
	layout := NewLayout("/repo", WithStateDirectory("/custom/state"))

	if layout.StateDirectory() != "/custom/state" {
		t.Fatalf("StateDirectory = %q, want %q", layout.StateDirectory(), "/custom/state")
	}

	if layout.ToolBinaryDirectory() != "/custom/state/bin" {
		t.Fatalf("ToolBinaryDirectory = %q, want %q",
			layout.ToolBinaryDirectory(), "/custom/state/bin")
	}
}
