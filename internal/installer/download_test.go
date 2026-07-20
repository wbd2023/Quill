package installer

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDownloadFileWritesContentToDisk(t *testing.T) {
	t.Parallel()

	const body = "archive-body"

	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		_ *http.Request,
	) {
		_, _ = io.WriteString(writer, body)
	}))
	defer server.Close()

	destination := filepath.Join(t.TempDir(), "downloads", "archive")
	if err := downloadFile(t.Context(), server.URL, destination); err != nil {
		t.Fatalf("download file: %v", err)
	}

	downloaded, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read download: %v", err)
	}
	if string(downloaded) != body {
		t.Fatalf("downloaded content = %q, want %q", downloaded, body)
	}
}

func TestDownloadFileRejectsOversizedResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(
		writer http.ResponseWriter,
		_ *http.Request,
	) {
		_, _ = io.WriteString(writer, "123456789")
	}))
	defer server.Close()

	directory := t.TempDir()
	destination := filepath.Join(directory, "archive")
	err := downloadFileUpTo(t.Context(), server.URL, destination, 8)
	if err == nil {
		t.Fatal("expected oversized download to fail")
	}
	if _, statErr := os.Stat(destination); !os.IsNotExist(statErr) {
		t.Fatalf("destination exists after failure: %v", statErr)
	}

	entries, readErr := os.ReadDir(directory)
	if readErr != nil {
		t.Fatalf("read download directory: %v", readErr)
	}
	if len(entries) != 0 {
		t.Fatalf("temporary downloads remain after failure: %v", entries)
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
