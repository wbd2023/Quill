package execution

import (
	"slices"
	"strings"
	"testing"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

/* ------------------------------------------ File Sets ----------------------------------------- */

func TestCollectFileSetFilesUsesProfileDefinedScopeFilters(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, scopedFileSetConfig(t))

	readme := testutil.WriteFile(t, repoRoot, "README.md", "# App\n")
	testutil.WriteFile(t, repoRoot, "internal/guide.md", "# Internal\n")
	toolsReadme := testutil.WriteFile(t, repoRoot, "tools/README.md", "# Tools\n")

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
	profiles.Write(t, repoRoot, scopedFileSetConfig(t))

	toolsReadme := testutil.WriteFile(t, repoRoot, "tools/README.md", "# Tools\n")

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
	profiles.Write(t, repoRoot, scopedFileSetConfig(t))

	testutil.WriteFile(t, repoRoot, "README.md", "# App\n")

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
	makefile := testutil.WriteFile(t, repoRoot, "Makefile", "all:\n\t@true\n")
	config := testutil.WriteFile(t, repoRoot, "style.local.toml", "enabled = true\n")
	checksum := testutil.WriteFile(t, repoRoot, "go.sum", strings.Repeat("x", 120)+"\n")

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

func TestCollectFileSetFilesSelectsPrivacyDocumentAndExcludesUnlistedRootMarkdown(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))

	privacy := testutil.WriteFile(t, repoRoot, "CONTRIBUTOR_PRIVACY.md", "# Contributor Privacy\n")
	unlisted := testutil.WriteFile(t, repoRoot, "UNLISTED.md", "# Unlisted\n")

	context := testContext(t, repoRoot, style.Scope("all"))

	files, err := CollectFileSetFiles(context, "markdown")
	if err != nil {
		t.Fatalf("collect markdown file set: %v", err)
	}

	if !slices.Contains(files, privacy) {
		t.Fatalf("expected CONTRIBUTOR_PRIVACY.md in markdown file set: %v", files)
	}

	if slices.Contains(files, unlisted) {
		t.Fatalf("did not expect unlisted root markdown in markdown file set: %v", files)
	}
}

func TestCollectLineLengthFileSetCoversPrivacyDocument(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Current(t))

	privacy := testutil.WriteFile(t, repoRoot, "CONTRIBUTOR_PRIVACY.md", "# Contributor Privacy\n")

	context := testContext(t, repoRoot, style.Scope("all"))

	files, err := CollectFileSetFiles(context, "line_length")
	if err != nil {
		t.Fatalf("collect line_length file set: %v", err)
	}

	if !slices.Contains(files, privacy) {
		t.Fatalf("expected CONTRIBUTOR_PRIVACY.md in line_length file set: %v", files)
	}
}

func scopedFileSetConfig(t *testing.T) (config policy.Config) {
	t.Helper()

	config = profiles.Current(t)
	config.Repository.ScopeRoots = map[style.Scope][]string{
		"app":   {"cmd", "internal", "test"},
		"tools": {"tools"},
		"all":   {"."},
	}
	config.Repository.DefaultScope = style.Scope("app")

	for index := range config.FileSets {
		if config.FileSets[index].Name != "markdown" {
			continue
		}

		config.FileSets[index].Include.Files = map[style.Scope][]string{
			"app": {"README.md"},
		}
		config.FileSets[index].Include.Paths = map[style.Scope][]string{
			"app":   {"internal/"},
			"tools": {"tools/"},
		}
		break
	}

	return config
}
