package installer

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type failingReader struct{}

func (failingReader) Read(_ []byte) (count int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestWriteExecutableReplacesRegularDestination(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	destination := filepath.Join(directory, "tool")
	if err := os.WriteFile(destination, []byte("previous"), 0o600); err != nil {
		t.Fatalf("write existing destination: %v", err)
	}

	if err := writeExecutable(
		directory,
		destination,
		strings.NewReader("replacement"),
	); err != nil {
		t.Fatalf("replace destination: %v", err)
	}

	content, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if string(content) != "replacement" {
		t.Fatalf("destination content = %q, want replacement", content)
	}

	info, err := os.Stat(destination)
	if err != nil {
		t.Fatalf("inspect destination: %v", err)
	}
	if info.Mode().Perm() != standardPermissions {
		t.Fatalf("destination mode = %v, want %v", info.Mode().Perm(), standardPermissions)
	}
}

func TestWriteExecutablePreservesDestinationAfterReadFailure(t *testing.T) {
	t.Parallel()

	directory := t.TempDir()
	destination := filepath.Join(directory, "tool")
	if err := os.WriteFile(destination, []byte("previous"), 0o600); err != nil {
		t.Fatalf("write existing destination: %v", err)
	}

	if err := writeExecutable(directory, destination, failingReader{}); err == nil {
		t.Fatal("expected reader failure")
	}

	content, err := os.ReadFile(destination)
	if err != nil {
		t.Fatalf("read destination: %v", err)
	}
	if string(content) != "previous" {
		t.Fatalf("destination content = %q, want previous", content)
	}

	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatalf("read destination directory: %v", err)
	}
	if len(entries) != 1 || entries[0].Name() != filepath.Base(destination) {
		t.Fatalf(
			"destination directory entries = %v, want only %q",
			entries,
			filepath.Base(destination),
		)
	}
}
