package lockfile

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

/* -------------------------------------- Lockfile Writing -------------------------------------- */

func TestWriteAppliesSharedFilePermissions(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	lockfile := Lockfile{} // empty archives is a valid (if useless) lockfile

	path, err := Write(root, lockfile)
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

	if info.Mode().Perm() != standardLockfilePermissions {
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

	path, err := Write(root, lockfile)
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

	path, err := Write(root, original)
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
