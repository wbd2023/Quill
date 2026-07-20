package cli

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteLockfileUsesSharedFilePermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "quill.lock")
	if err := writeLockfile(path, "schema_version = 1\n"); err != nil {
		t.Fatalf("writeLockfile: %v", err)
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
