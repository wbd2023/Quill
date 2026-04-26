package runtime

import (
	"archive/tar"
	"os"
	"path/filepath"
	"testing"
)

func TestExtractShellcheckBinaryExtractsExpectedAsset(t *testing.T) {
	t.Parallel()

	archivePath := writeShellcheckArchive(
		t,
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/LICENSE.txt",
			Body: "licence",
		},
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/README.txt",
			Body: "readme",
		},
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/shellcheck",
			Body: "#!/bin/sh\n",
		},
	)

	destination := t.TempDir()
	binaryPath, err := extractShellcheckBinary(archivePath, destination, "0.10.0")
	if err != nil {
		t.Fatalf("extractShellcheckBinary: %v", err)
	}

	contents, err := os.ReadFile(binaryPath)
	if err != nil {
		t.Fatalf("read extracted binary: %v", err)
	}

	if string(contents) != "#!/bin/sh\n" {
		t.Fatalf("unexpected extracted contents: %q", contents)
	}

	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Fatalf("stat extracted binary: %v", err)
	}

	if info.Mode().Perm() != executableMode {
		t.Fatalf("extracted binary mode = %v, want %v", info.Mode().Perm(), executableMode)
	}

	readmePath := filepath.Join(destination, "shellcheck-v0.10.0", "README.txt")
	if _, err := os.Stat(readmePath); err == nil {
		t.Fatal("expected README to be ignored, not extracted")
	}
}

func TestExtractShellcheckBinaryRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	archivePath := writeShellcheckArchive(
		t,
		shellcheckArchiveEntry{
			Name: "shellcheck-v0.10.0/../shellcheck",
			Body: "bad",
		},
	)

	if _, err := extractShellcheckBinary(archivePath, t.TempDir(), "0.10.0"); err == nil {
		t.Fatal("expected path traversal entry to fail")
	}
}

func TestExtractShellcheckBinaryRejectsLinks(t *testing.T) {
	t.Parallel()

	archivePath := writeShellcheckArchive(
		t,
		shellcheckArchiveEntry{
			Name:     "shellcheck-v0.10.0/shellcheck",
			Typeflag: tar.TypeSymlink,
			Linkname: "/bin/sh",
		},
	)

	if _, err := extractShellcheckBinary(archivePath, t.TempDir(), "0.10.0"); err == nil {
		t.Fatal("expected link entry to fail")
	}
}
