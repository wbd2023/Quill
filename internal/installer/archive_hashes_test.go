package installer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestArchiveHashForReturnsKnownPlatform(t *testing.T) {
	t.Parallel()

	hash, err := archiveHashFor("shellcheck", "linux", "amd64")
	if err != nil {
		t.Fatalf("archiveHashFor: %v", err)
	}

	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
}

func TestArchiveHashForRejectsUnknownPlatform(t *testing.T) {
	t.Parallel()

	if _, err := archiveHashFor("shellcheck", "freebsd", "amd64"); err == nil {
		t.Fatal("expected unsupported platform error")
	}
}

func TestArchiveHashForRejectsUnknownTool(t *testing.T) {
	t.Parallel()

	if _, err := archiveHashFor("nonexistent", "linux", "amd64"); err == nil {
		t.Fatal("expected unknown tool error")
	}
}

func TestVerifyFileChecksumRejectsMismatch(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "archive.tar.xz")
	if err := os.WriteFile(path, []byte("archive"), 0o600); err != nil {
		t.Fatalf("write archive: %v", err)
	}

	err := verifyChecksum(path, strings.Repeat("0", 64))
	if err == nil {
		t.Fatal("expected checksum mismatch")
	}
}
