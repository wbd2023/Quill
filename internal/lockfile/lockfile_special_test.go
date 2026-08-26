//go:build unix

package lockfile

import (
	"net"
	"path/filepath"
	"strings"
	"testing"
)

// TestLoadRejectsNonRegularLockfile guards Load against special-file leaves. A Unix domain
// socket at the lockfile path must fail closed at the regular-file check, before Load opens
// or reads anything.
func TestLoadRejectsNonRegularLockfile(t *testing.T) {
	t.Parallel()

	root := t.TempDir()
	socket := filepath.Join(root, DefaultFilename)
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Skipf("unix socket unsupported: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() }) // the listener exists only to create the socket path

	lockfile, err := Load(root)
	if err == nil {
		t.Fatal("expected error for non-regular lockfile, got nil")
	}
	if lockfile.Loaded {
		t.Fatalf("non-regular lockfile must not load: %+v", lockfile)
	}
	if !strings.Contains(err.Error(), "not a regular file") {
		t.Fatalf("error %q does not mention non-regular file", err.Error())
	}
}
