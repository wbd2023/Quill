package lockfile

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

/* -------------------------------------- Lockfile Writing -------------------------------------- */

func TestWriteAppliesSharedFilePermissions(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	lockfile := Lockfile{} // empty archives is a valid (if useless) lockfile

	path, err := Write(context.Background(), root, lockfile)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	if want := filepath.Join(root, DefaultFilename); path != want {
		t.Fatalf("path = %q, want %q", path, want)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat lockfile: %v", err)
	}

	if runtime.GOOS != "windows" && info.Mode().Perm() != standardLockfilePermissions {
		t.Fatalf(
			"lockfile permissions = %04o, want %04o",
			info.Mode().Perm(),
			standardLockfilePermissions,
		)
	}
}

func TestWriteCreatesParentDirectories(t *testing.T) {
	t.Parallel()

	// Write writes to <root>/quill.lock; the root itself may not exist yet
	// when a caller has only resolved a logical root.
	root := filepath.Join(t.TempDir(), "nested", "repo")
	lockfile := Lockfile{
		Archives: map[string]Archive{
			"shellcheck": {
				Tool:    "shellcheck",
				Version: "0.10.0",
				Hashes:  map[string]string{"linux/amd64": "abc"},
			},
		},
	}

	path, err := Write(context.Background(), root, lockfile)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read written lockfile: %v", err)
	}

	if !strings.Contains(string(contents), "shellcheck") {
		t.Fatalf("written lockfile missing shellcheck entry:\n%s", string(contents))
	}
}

func TestWriteRoundTripsThroughDecode(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	original := Lockfile{
		Archives: map[string]Archive{
			"shellcheck": {
				Tool:    "shellcheck",
				Version: "0.10.0",
				Hashes:  map[string]string{"linux/amd64": "abc", "darwin/arm64": "def"},
			},
		},
	}

	path, err := Write(context.Background(), root, original)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	loaded, err := Load(root)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	archive, ok := loaded.Archives["shellcheck"]
	if !ok {
		t.Fatalf("loaded lockfile missing shellcheck after Write to %s", path)
	}

	if archive.Version != "0.10.0" {
		t.Fatalf("version = %q, want 0.10.0", archive.Version)
	}

	if len(archive.Hashes) != 2 {
		t.Fatalf("expected 2 hashes, got %d", len(archive.Hashes))
	}
}

func TestWriteDoesNotReplaceAfterCancellationAtCommit(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	path := filepath.Join(root, DefaultFilename)
	original := []byte("preserve existing lockfile\n")
	if err := os.WriteFile(path, original, standardLockfilePermissions); err != nil {
		t.Fatal(err)
	}

	ctx := &cancellationAtCommitContext{done: make(chan struct{})}
	_, err := Write(ctx, root, Lockfile{})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Write error = %v, want context.Canceled", err)
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != string(original) {
		t.Fatalf("lockfile contents = %q, want %q", contents, original)
	}
}

type cancellationAtCommitContext struct {
	done   chan struct{}
	checks int
}

func (*cancellationAtCommitContext) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (ctx *cancellationAtCommitContext) Done() (done <-chan struct{}) {
	return ctx.done
}

func (ctx *cancellationAtCommitContext) Err() (err error) {
	ctx.checks++
	if ctx.checks < 2 {
		return nil
	}
	close(ctx.done)
	return context.Canceled
}

func (*cancellationAtCommitContext) Value(any) (value any) {
	return nil
}
