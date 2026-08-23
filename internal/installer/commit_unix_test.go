//go:build !windows

package installer

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWriteExecutableResistsParentSwapDuringWrite proves the rooted commit cannot be redirected
// outside the repository when the destination's parent directory is swapped to a symlink mid-write.
//
// A blocking reader displaces the anchored parent directory and replaces its path with a symlink to
// a directory outside the repository while the executable body is being staged. A name-based commit
// would resolve the destination through that symlink; the rooted commit must not.
func TestWriteExecutableResistsParentSwapDuringWrite(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	outside := t.TempDir()

	binDir := filepath.Join(root, ".cache", "quill", "bin")
	if err := os.MkdirAll(binDir, standardPermissions); err != nil {
		t.Fatalf("create bin directory: %v", err)
	}

	// Hostile sentinel living outside the repository root. A redirected commit would overwrite it.
	sentinelDir := filepath.Join(outside, "quill", "bin")
	if err := os.MkdirAll(sentinelDir, standardPermissions); err != nil {
		t.Fatalf("create sentinel directory: %v", err)
	}
	sentinel := filepath.Join(sentinelDir, "tool")
	if err := os.WriteFile(sentinel, []byte("protected"), 0o600); err != nil {
		t.Fatalf("write sentinel: %v", err)
	}

	destination := filepath.Join(binDir, "tool")

	reader := &swappingReader{
		swap: func() error {
			// Displace the anchored parent and replace its path with a symlink pointing outside the
			// repository, so a name-based commit would resolve the destination through the symlink.
			displaced := binDir + ".displaced"
			if err := os.Rename(binDir, displaced); err != nil {
				return err
			}
			return os.Symlink(sentinelDir, binDir)
		},
		data: []byte("replacement"),
	}

	writeErr := writeExecutable(root, destination, reader)

	// The commit must never reach the outside sentinel, regardless of the outcome.
	content, err := os.ReadFile(sentinel)
	if err != nil {
		t.Fatalf("read sentinel: %v", err)
	}
	if string(content) != "protected" {
		t.Fatalf("outside sentinel was overwritten: %q", content)
	}

	// The rooted commit must succeed against the anchored directory despite the swap, landing the
	// executable inside the repository rather than failing closed or escaping.
	if writeErr != nil {
		t.Fatalf("writeExecutable failed under parent swap: %v", writeErr)
	}

	written := filepath.Join(binDir+".displaced", "tool")
	got, err := os.ReadFile(written)
	if err != nil {
		t.Fatalf("read written executable: %v", err)
	}
	if string(got) != "replacement" {
		t.Fatalf("written executable content = %q, want %q", got, "replacement")
	}

	info, statErr := os.Stat(written)
	if statErr != nil {
		t.Fatalf("stat written executable: %v", statErr)
	}
	if info.Mode().Perm() != standardPermissions {
		t.Fatalf("written executable mode = %v, want %v", info.Mode().Perm(), standardPermissions)
	}
}
