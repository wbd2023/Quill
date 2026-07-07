package installer

import (
	"archive/tar"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractShellcheckBinaryExtractsExpectedAsset(t *testing.T) {
	t.Parallel()

	archive := writeShellcheckArchive(
		t,
		archiveEntry{
			Name: "shellcheck-v0.10.0/LICENSE.txt",
			Body: "licence",
		},
		archiveEntry{
			Name: "shellcheck-v0.10.0/README.txt",
			Body: "readme",
		},
		archiveEntry{
			Name: "shellcheck-v0.10.0/shellcheck",
			Body: "#!/bin/sh\n",
		},
	)

	dir := t.TempDir()
	binary, err := extractShellcheckBinary(archive, dir, "0.10.0")
	if err != nil {
		t.Fatalf("extract: %v", err)
	}

	content, err := os.ReadFile(binary)
	if err != nil {
		t.Fatalf("read extracted binary: %v", err)
	}

	if string(content) != "#!/bin/sh\n" {
		t.Fatalf("unexpected extracted content: %q", content)
	}

	info, err := os.Stat(binary)
	if err != nil {
		t.Fatalf("stat extracted binary: %v", err)
	}

	if info.Mode().Perm() != standardPermissions {
		t.Fatalf("extracted binary mode = %v, want %v", info.Mode().Perm(), standardPermissions)
	}

	readme := filepath.Join(dir, "shellcheck-v0.10.0", "README.txt")
	if _, err := os.Stat(readme); err == nil {
		t.Fatal("expected README to be ignored, not extracted")
	}
}

func TestExtractShellcheckBinaryRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	archive := writeShellcheckArchive(
		t,
		archiveEntry{
			Name: "shellcheck-v0.10.0/../shellcheck",
			Body: "bad",
		},
	)

	if _, err := extractShellcheckBinary(archive, t.TempDir(), "0.10.0"); err == nil {
		t.Fatal("expected path traversal entry to fail")
	}
}

func TestExtractShellcheckBinaryRejectsLinks(t *testing.T) {
	t.Parallel()

	archive := writeShellcheckArchive(
		t,
		archiveEntry{
			Name:     "shellcheck-v0.10.0/shellcheck",
			Typeflag: tar.TypeSymlink,
			Linkname: "/bin/sh",
		},
	)

	if _, err := extractShellcheckBinary(archive, t.TempDir(), "0.10.0"); err == nil {
		t.Fatal("expected link entry to fail")
	}
}
