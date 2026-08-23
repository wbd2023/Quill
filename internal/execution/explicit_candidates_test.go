package execution

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
)

// TestExplicitFileCandidatesExcludesSymlinkLeaf pins the fail-closed contract for explicit
// include.files entries: a symlink leaf is dropped before it reaches a driver, even though it
// resolves to a regular file. The historical os.Stat followed the link and admitted it.
func TestExplicitFileCandidatesExcludesSymlinkLeaf(t *testing.T) {
	repoRoot := t.TempDir()
	regular := testutil.WriteFile(t, repoRoot, "kept.sh", "echo hi\n")

	target := testutil.WriteFile(t, t.TempDir(), "secret.sh", "echo secret\n")
	evil := filepath.Join(repoRoot, "evil.sh")
	if err := os.Symlink(target, evil); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	context := RunContext{RepoRoot: repoRoot}
	fileSet := profile.FileSetConfig{
		Include: profile.FileSetInclude{
			Files: map[style.Scope][]string{
				style.Scope("all"): {"kept.sh", "evil.sh"},
			},
		},
	}

	files := explicitFileCandidates(context, fileSet, []style.Scope{style.Scope("all")})

	if !slices.Contains(files, regular) {
		t.Fatalf("expected regular explicit candidate %q in %v", regular, files)
	}
	if slices.Contains(files, evil) {
		t.Fatalf("symlink explicit candidate %q must be excluded: %v", evil, files)
	}
}
