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
}

func TestSearchPathPrependsBinaryDirectories(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")

	actual := SearchPath("/tool/bin", "/node/bin")
	expected := "/tool/bin:/node/bin:/usr/bin"
	if actual != expected {
		t.Fatalf("SearchPath = %q, want %q", actual, expected)
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
