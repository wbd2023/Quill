//go:build unix

package filewalk

import (
	"net"
	"path/filepath"
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

// TestCollectFilesExcludesSpecialFileLeafMatchingExtension guards the walk gate against
// non-regular leaves other than symlinks. A Unix domain socket named with a matching
// extension must be excluded before any probe opens it; the old extension-only matcher
// would have collected it.
func TestCollectFilesExcludesSpecialFileLeafMatchingExtension(t *testing.T) {
	root := t.TempDir()
	regular := testutil.WriteFile(t, root, "kept.sh", "#!/bin/sh\necho hi\n")

	socket := filepath.Join(root, "blocked.sh")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Skipf("unix socket unsupported: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() }) // the listener exists only to create the socket path

	files, err := CollectFiles([]string{root}, walkConfig(profiles.RepositoryConfig()), ".sh")
	if err != nil {
		t.Fatalf("CollectFiles: %v", err)
	}

	requirePaths(t, files, []string{regular})
}
