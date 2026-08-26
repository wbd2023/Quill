//go:build windows

package installer

import (
	"os"
	"path/filepath"
	"testing"
)

// TestWriteExecutableResistsParentSwapDuringWrite proves the rooted Windows commit cannot be
// redirected outside the repository when the destination's parent is swapped to a symlink
// mid-write.
//
// The NT handle-relative rename resolves the staged file and its destination against the anchored
// directory descriptor, so the displaced parent path (and the symlink that replaces it) cannot move
// the destination outside the repository. Creating the redirecting symlink needs
// SeCreateSymbolicLinkPrivilege, so the test skips where that privilege is absent.
func TestWriteExecutableResistsParentSwapDuringWrite(t *testing.T) {
	// Probe symlink privilege before staging the hostile layout.
	probeDir := t.TempDir()
	if err := os.Symlink(probeDir, filepath.Join(probeDir, "probe")); err != nil {
		t.Skipf("parent-swap regression requires symlink privilege: %v", err)
	}

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

	// The rooted commit must land the executable inside the anchored repository tree.
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
	if info.Mode().Perm()&0o200 == 0 {
		t.Fatalf("written executable mode = %v, want writable file", info.Mode().Perm())
	}
}
