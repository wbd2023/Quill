package runner

import (
	"slices"
	"strings"
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/style"
)

/* ------------------------------------------ File Sets ----------------------------------------- */

func TestCollectFileSetFilesUsesProfileDefinedScopeFilters(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))

	readme := fixtures.WriteFile(t, repoRoot, "README.md", "# App\n")
	fixtures.WriteFile(t, repoRoot, "internal/guide.md", "# Internal\n")
	toolsReadme := fixtures.WriteFile(t, repoRoot, "tools/README.md", "# Tools\n")

	context := testContext(t, repoRoot, style.Scope("app"))

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

func TestCollectFileSetFilesDoesNotUseDefaultScopeAsWidestScope(t *testing.T) {
	repoRoot := t.TempDir()
	config := profiles.Current(t)
	config.Repository.DefaultScope = style.Scope("app")
	profiles.Write(t, repoRoot, config)

	toolsReadme := fixtures.WriteFile(t, repoRoot, "tools/README.md", "# Tools\n")

	context := testContext(t, repoRoot, style.Scope("all"))

	files, err := CollectFileSetFiles(context, "markdown")
	if err != nil {
		t.Fatalf("collect markdown file set: %v", err)
	}

	if !slices.Contains(files, toolsReadme) {
		t.Fatalf("expected tools markdown file despite app default scope: %v", files)
	}
}

func TestCollectFileSetFilesReturnsEmptySetWhenScopedIncludesDoNotOverlap(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))

	fixtures.WriteFile(t, repoRoot, "README.md", "# App\n")

	context := testContext(t, repoRoot, style.Scope("tools"))

	files, err := CollectFileSetFiles(context, "markdown")
	if err != nil {
		t.Fatalf("collect markdown file set: %v", err)
	}

	if len(files) != 0 {
		t.Fatalf("expected no markdown files in tools scope: %v", files)
	}
}

func TestCollectFileSetFilesRejectsUnknownSet(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))

	context := testContext(t, repoRoot, style.Scope("all"))

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

	context := testContext(t, repoRoot, style.Scope("all"))
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
