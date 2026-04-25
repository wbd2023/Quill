package runner

import (
	"slices"
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

/* ------------------------------------------ File Sets ----------------------------------------- */

func TestCollectFileSetFilesUsesProfileDefinedScopeFilters(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))

	readme := fixtures.WriteFile(t, repoRoot, "README.md", "# App\n")
	fixtures.WriteFile(t, repoRoot, "internal/guide.md", "# Internal\n")
	toolsReadme := fixtures.WriteFile(t, repoRoot, "tools/README.md", "# Tools\n")

	context := testContext(t, repoRoot, contract.ScopeApp)

	files, err := CollectFileSetFiles(context, "markdown")
	if err != nil {
		t.Fatalf("collect markdown file set: %v", err)
	}

	if !slices.Contains(files, readme) {
		t.Fatalf("expected app markdown file in %v", files)
	}

	if slices.Contains(files, toolsReadme) {
		t.Fatalf("did not expect tools markdown file in app scope: %v", files)
	}
}

func TestCollectFileSetFilesRejectsUnknownSet(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))

	context := testContext(t, repoRoot, contract.ScopeAll)

	if _, err := CollectFileSetFiles(context, "javascript"); err == nil {
		t.Fatal("expected unknown file set to fail")
	}
}

func TestCollectLineLengthFileSetCoversTextFiles(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))
	makefile := fixtures.WriteFile(t, repoRoot, "Makefile", "all:\n\t@true\n")
	config := fixtures.WriteFile(t, repoRoot, "style.local.toml", "enabled = true\n")
	checksum := fixtures.WriteFile(t, repoRoot, "go.sum", strings.Repeat("x", 120)+"\n")

	context := testContext(t, repoRoot, contract.ScopeAll)
	files, err := CollectFileSetFiles(context, "line_length")
	if err != nil {
		t.Fatalf("collect line_length file set: %v", err)
	}

	if !slices.Contains(files, makefile) {
		t.Fatalf("expected Makefile in line_length file set: %v", files)
	}

	if !slices.Contains(files, config) {
		t.Fatalf("expected TOML config in line_length file set: %v", files)
	}

	if slices.Contains(files, checksum) {
		t.Fatalf("expected go.sum to be excluded from line_length file set: %v", files)
	}
}
