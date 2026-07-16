package workspace

import "testing"

func TestLayoutDerivesPathsFromRepositoryRoot(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")

	layout := NewLayout("/repo")

	if layout.StateDirectory != "/repo/.cache/quill" {
		t.Fatalf("StateDirectory = %q, want %q", layout.StateDirectory, "/repo/.cache/quill")
	}

	if layout.BinaryDirectory() != "/repo/.cache/quill/bin" {
		t.Fatalf("BinaryDirectory = %q, want %q",
			layout.BinaryDirectory(), "/repo/.cache/quill/bin")
	}
}

func TestLayoutBuildPathIncludesBinaryDirectoryAndExtras(t *testing.T) {
	t.Setenv("PATH", "/usr/bin")

	layout := NewLayout("/repo")

	actual := layout.BuildPath("/node/bin")
	expected := "/repo/.cache/quill/bin:/node/bin:/usr/bin"
	if actual != expected {
		t.Fatalf("BuildPath = %q, want %q", actual, expected)
	}
}

func TestLayoutStateDirectoryCanBeOverridden(t *testing.T) {
	layout := Layout{
		RepositoryRoot: "/repo",
		StateDirectory: "/custom/state",
	}

	if layout.BinaryDirectory() != "/custom/state/bin" {
		t.Fatalf("BinaryDirectory = %q, want %q",
			layout.BinaryDirectory(), "/custom/state/bin")
	}
}
