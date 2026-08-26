package cli

import (
	"testing"

	"github.com/wbd2023/quill/internal/testutil"
)

func TestResolveRepositoryRootAutoDetectsRepository(t *testing.T) {
	root, err := resolveRepositoryRoot("")
	if err != nil {
		t.Fatalf("resolveRepositoryRoot: %v", err)
	}

	if root != testutil.RepositoryRoot(t) {
		t.Fatalf("unexpected repo root %q", root)
	}
}

func TestFindRepoRootRejectsMissingRepository(t *testing.T) {
	missingRoot := t.TempDir()

	_, err := findRepositoryRoot(missingRoot)
	if err == nil {
		t.Fatal("expected missing repository error")
	}
}

func TestFindRepoRootRejectsLegacyRootWithoutQuillProfile(t *testing.T) {
	legacyRoot := t.TempDir()
	testutil.WriteFile(t, legacyRoot, "STYLE.md", "# Style\n")

	testutil.WriteFile(t, legacyRoot, "tools/go.mod", "module example.com/tools\n")

	if _, err := findRepositoryRoot(legacyRoot); err == nil {
		t.Fatal("expected legacy repo root without quill.toml to be rejected")
	}
}
