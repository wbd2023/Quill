package engine

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/wbd2023/quill/internal/lockfile"
)

func TestWriteResolvedLockDoesNotPublishAfterCancellation(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, lockfile.DefaultFilename)
	original := []byte("preserve this lockfile\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := writeResolvedLock(ctx, root, []lockfile.Archive{{
		Tool:    "tool",
		Version: "1.0.0",
		Hashes:  map[string]string{"linux/amd64": "checksum"},
	}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("writeResolvedLock error = %v, want context.Canceled", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != string(original) {
		t.Fatalf("lockfile contents = %q, want %q", contents, original)
	}
}
