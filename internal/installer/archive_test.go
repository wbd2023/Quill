package installer

import (
	"archive/tar"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/wbd2023/Quill/internal/toolchain"
)

// testInstall mirrors the shellcheck GitHubInstall from
// pack/shipped/tool/builders.go without importing that package, so the
// installer's archive extraction is tested against a fixed shape.
func testInstall() (install toolchain.GitHubInstall) {
	return toolchain.GitHubInstall{
		Tag:  "v%s",
		Path: "shellcheck-%s/shellcheck",
	}
}

/* -------------------------------------------- Tests ------------------------------------------- */

func TestExtractBinaryExtractsExpectedAsset(t *testing.T) {
	t.Parallel()

	archive := writeTestArchive(
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
	binary, err := extractBinary(archive, dir, testInstall(), "0.10.0")
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

func TestExtractBinaryRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	archive := writeTestArchive(
		t,
		archiveEntry{
			Name: "shellcheck-v0.10.0/../shellcheck",
			Body: "bad",
		},
	)

	if _, err := extractBinary(archive, t.TempDir(), testInstall(), "0.10.0"); err == nil {
		t.Fatal("expected path traversal entry to fail")
	}
}

func TestExtractBinaryRejectsLinks(t *testing.T) {
	t.Parallel()

	archive := writeTestArchive(
		t,
		archiveEntry{
			Name:     "shellcheck-v0.10.0/shellcheck",
			Typeflag: tar.TypeSymlink,
			Linkname: "/bin/sh",
		},
	)

	if _, err := extractBinary(archive, t.TempDir(), testInstall(), "0.10.0"); err == nil {
		t.Fatal("expected link entry to fail")
	}
}

func TestExtractBinaryRejectsOversizedEntry(t *testing.T) {
	t.Parallel()

	const oversizedEntry = int64(maxArchiveSize) + 1

	archive := writeTestArchiveHeader(
		t,
		"shellcheck-v0.10.0/shellcheck",
		oversizedEntry,
	)
	dir := t.TempDir()
	_, err := extractBinary(archive, dir, testInstall(), "0.10.0")
	if err == nil || !strings.Contains(err.Error(), "uncompressed size") {
		t.Fatalf("extract oversized entry error = %v, want uncompressed size error", err)
	}

	target := filepath.Join(dir, "shellcheck-v0.10.0", "shellcheck")
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("target exists after oversized archive: %v", statErr)
	}
}

func TestExtractBinaryRejectsCumulativeUncompressedSize(t *testing.T) {
	t.Parallel()

	archive := writeTestArchive(
		t,
		archiveEntry{
			Name: "shellcheck-v0.10.0/LICENSE.txt",
			Body: "12345",
		},
		archiveEntry{
			Name: "shellcheck-v0.10.0/shellcheck",
			Body: "6789",
		},
	)
	dir := t.TempDir()
	_, err := extractBinaryUpTo(archive, dir, testInstall(), "0.10.0", 8)
	if err == nil || !strings.Contains(err.Error(), "uncompressed size") {
		t.Fatalf("extract cumulative size error = %v, want uncompressed size error", err)
	}

	target := filepath.Join(dir, "shellcheck-v0.10.0", "shellcheck")
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Fatalf("target exists after oversized archive: %v", statErr)
	}
}
