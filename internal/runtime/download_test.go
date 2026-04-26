package runtime

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteDownloadStreamsBodyToDisk(t *testing.T) {
	t.Parallel()

	body := "archive-body"

	destination := filepath.Join(t.TempDir(), "archive.tar.xz")
	if err := writeDownload(destination, strings.NewReader(body)); err != nil {
		t.Fatalf("writeDownload: %v", err)
	}

	contents, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read downloaded file: %v", err)
	}

	if string(contents) != body {
		t.Fatalf("downloaded contents = %q", contents)
	}

	info, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("stat downloaded file: %v", err)
	}

	if info.Mode().Perm() != downloadMode {
		t.Fatalf("downloaded file mode = %v, want %v", info.Mode().Perm(), downloadMode)
	}

	expectedChecksum := fmt.Sprintf("%x", sha256.Sum256([]byte(body)))
	if err := verifyFileChecksum(destination, "archive.tar.xz", expectedChecksum); err != nil {
		t.Fatalf("verify downloaded checksum: %v", err)
	}
}

func TestWriteDownloadRejectsOversizedBody(t *testing.T) {
	t.Parallel()

	destination := filepath.Join(t.TempDir(), "archive.tar.xz")
	err := writeDownloadWithLimit(destination, io.LimitReader(zeroReader{}, 5), 4)
	if err == nil || !strings.Contains(err.Error(), "maximum size") {
		t.Fatalf("expected maximum size error, got %v", err)
	}
}

type zeroReader struct{}

func (zeroReader) Read(buffer []byte) (count int, err error) {
	for index := range buffer {
		buffer[index] = 0
	}

	return len(buffer), nil
}
