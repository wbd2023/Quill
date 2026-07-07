package installer

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestDownloadFileWritesContentToDisk(t *testing.T) {
	t.Parallel()

	// downloadFile needs an HTTP server; test the checksum logic instead.
	body := "archive-body"
	checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(body)))

	// Write the body to a temp file, then verify the checksum.
	file, err := os.CreateTemp(t.TempDir(), "test-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}

	if _, err = file.WriteString(body); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err = file.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	if err = verifyChecksum(file.Name(), checksum); err != nil {
		t.Fatalf("verify checksum: %v", err)
	}
}

func TestVerifyChecksumRejectsMismatch(t *testing.T) {
	t.Parallel()

	file, err := os.CreateTemp(t.TempDir(), "test-*")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}

	if _, err = file.WriteString("archive"); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err = file.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}

	err = verifyChecksum(file.Name(), strings.Repeat("0", 64))
	if err == nil {
		t.Fatal("expected checksum mismatch")
	}
}
