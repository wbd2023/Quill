//go:build unix

package execution

import (
	"net"
	"path/filepath"
	"slices"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
)

// TestExplicitFileCandidatesExcludesSpecialFileLeaf confirms non-regular leaves beyond
// symlinks are rejected by the explicit-file path. A Unix domain socket is never a
// candidate, regardless of the name configured in include.files.
func TestExplicitFileCandidatesExcludesSpecialFileLeaf(t *testing.T) {
	repoRoot := t.TempDir()
	regular := testutil.WriteFile(t, repoRoot, "kept.sh", "echo hi\n")

	socket := filepath.Join(repoRoot, "blocked.sh")
	listener, err := net.Listen("unix", socket)
	if err != nil {
		t.Skipf("unix socket unsupported: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })

	context := RunContext{RepoRoot: repoRoot}
	fileSet := profile.FileSetConfig{
		Include: profile.FileSetInclude{
			Files: map[style.Scope][]string{
				style.Scope("all"): {"kept.sh", "blocked.sh"},
			},
		},
	}

	files := explicitFileCandidates(context, fileSet, []style.Scope{style.Scope("all")})

	if !slices.Contains(files, regular) {
		t.Fatalf("expected regular explicit candidate %q in %v", regular, files)
	}
	if slices.Contains(files, socket) {
		t.Fatalf("socket explicit candidate %q must be excluded: %v", socket, files)
	}
}
