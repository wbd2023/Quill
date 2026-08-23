package external_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/external"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

/* ------------------------------------------- Loading ------------------------------------------ */

func TestLoadSourcesProducesDefinition(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	packDir := filepath.Join(repoRoot, ".quill", "packs", "company")
	binDir := filepath.Join(packDir, "bin")
	mkdirAll(t, binDir)
	writeFile(t, filepath.Join(packDir, "pack.toml"), validManifest)
	writeExecutable(t, filepath.Join(binDir, "company-quill"))

	definitions, err := external.LoadSources(repoRoot, []profile.PackSource{
		{Path: ".quill/packs/company"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(definitions) != 1 {
		t.Fatalf("expected 1 definition, got %d", len(definitions))
	}

	definition := definitions[0]
	if definition.ID != "company" || definition.Name != "Company Engineering Policy" {
		t.Fatalf("unexpected definition: %+v", definition)
	}
	if len(definition.Rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(definition.Rules))
	}

	template, ok := definition.Rules[0].Check.(style.ExternalCheck)
	if !ok {
		t.Fatalf("expected ExternalCheck, got %T", definition.Rules[0].Check)
	}
	if template.CheckID != "no-direct-database-access" {
		t.Fatalf("unexpected check id: %q", template.CheckID)
	}
	if template.PackDirectory == "" {
		t.Fatal("expected resolved pack directory on the template")
	}
}

func TestLoadSourcesDefaultsOmittedRuleNameToRuleID(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	packDir := filepath.Join(repoRoot, ".quill", "packs", "company")
	binDir := filepath.Join(packDir, "bin")
	mkdirAll(t, binDir)
	manifest := strings.Replace(
		validManifest,
		`name = "No direct database access"`+"\n",
		"",
		1,
	)
	writeFile(t, filepath.Join(packDir, "pack.toml"), manifest)
	writeExecutable(t, filepath.Join(binDir, "company-quill"))

	definitions, err := external.LoadSources(repoRoot, []profile.PackSource{
		{Path: ".quill/packs/company"},
	})
	if err != nil {
		t.Fatalf("LoadSources: %v", err)
	}
	if got, want := definitions[0].Rules[0].Name, "company/no-direct-database-access"; got != want {
		t.Fatalf("Rule name = %q, want %q", got, want)
	}
}

func TestLoadSourcesRejectsMissingManifest(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	mkdirAll(t, filepath.Join(repoRoot, "empty"))

	if _, err := external.LoadSources(repoRoot, []profile.PackSource{{Path: "empty"}}); err == nil {
		t.Fatal("expected missing manifest error")
	}
}

func TestLoadSourcesRejectsUnresolvableExecutable(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	packDir := filepath.Join(repoRoot, ".quill", "packs", "company")
	mkdirAll(t, packDir)
	writeFile(t, filepath.Join(packDir, "pack.toml"), validManifest)

	if _, err := external.LoadSources(repoRoot, []profile.PackSource{
		{Path: ".quill/packs/company"},
	}); err == nil {
		t.Fatal("expected missing executable error before any launch")
	}
}

func TestLoadSourcesRejectsEscapingPath(t *testing.T) {
	t.Parallel()

	repoRoot := t.TempDir()
	outside := t.TempDir()
	packDir := filepath.Join(outside, "escaped")
	mkdirAll(t, packDir)
	writeFile(t, filepath.Join(packDir, "pack.toml"), validManifest)
	mkdirAll(t, filepath.Join(packDir, "bin"))
	writeExecutable(t, filepath.Join(packDir, "bin", "company-quill"))

	relative, err := filepath.Rel(repoRoot, packDir)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := external.LoadSources(repoRoot, []profile.PackSource{{Path: relative}}); err == nil {
		t.Fatal("expected repository escape rejection")
	}
}

func TestLoadSourcesRejectsSymlinkedManifestEscape(t *testing.T) {
	t.Parallel()
	if runtime.GOOS == "windows" {
		t.Skip("symlink semantics differ on Windows")
	}

	repoRoot := t.TempDir()
	packDir := filepath.Join(repoRoot, ".quill", "packs", "company")
	mkdirAll(t, packDir)
	mkdirAll(t, filepath.Join(packDir, "bin"))
	writeExecutable(t, filepath.Join(packDir, "bin", "company-quill"))

	outside := filepath.Join(t.TempDir(), "pack.toml")
	writeFile(t, outside, validManifest)

	if err := os.Symlink(outside, filepath.Join(packDir, "pack.toml")); err != nil {
		t.Skipf("symlink unsupported: %v", err)
	}

	if _, err := external.LoadSources(repoRoot, []profile.PackSource{
		{Path: ".quill/packs/company"},
	}); err == nil {
		t.Fatal("expected symlinked manifest escape rejection")
	}
}

/* ------------------------------------ Catalogue Composition ----------------------------------- */

// TestDuplicateExternalPackIDConflicts proves that two external Packs sharing an ID fail at
// catalogue assembly, the same boundary Shipped Packs use, so no external process runs.
func TestDuplicateExternalPackIDConflicts(t *testing.T) {
	t.Parallel()

	duplicate := pack.Definition{ID: "dup", Name: "Duplicate"}
	other := pack.Definition{ID: "dup", Name: "Other"}

	if _, err := shipped.ComposeCatalog(
		[]pack.Definition{duplicate, other},
	).Registry(nil); err == nil {
		t.Fatal("expected duplicate pack id conflict")
	}
}

// TestDuplicateRuleIDAcrossSourcesRejectedEvenWhenDisabled proves rule IDs are validated globally
// across every catalogue Pack before selection: two external Packs sharing a Rule ID conflict
// even when only one is enabled, so a disabled Pack cannot smuggle a colliding Rule ID into the
// catalogue.
func TestDuplicateRuleIDAcrossSourcesRejectedEvenWhenDisabled(t *testing.T) {
	t.Parallel()

	enabled := pack.Definition{
		ID:   "alpha",
		Name: "Alpha",
		Rules: []style.RuleDefinition{{
			ID:    "shared/rule",
			Check: style.ExternalCheck{CheckID: "c"},
		}},
	}
	disabled := pack.Definition{
		ID:   "beta",
		Name: "Beta",
		Rules: []style.RuleDefinition{{
			ID:    "shared/rule",
			Check: style.ExternalCheck{CheckID: "c"},
		}},
	}

	if _, err := shipped.ComposeCatalog(
		[]pack.Definition{enabled, disabled},
	).Registry([]string{"alpha"}); err == nil {
		t.Fatal("expected duplicate rule id rejection even with one pack disabled")
	}
}

// TestExternalPackComposesWithShipped proves an external Pack with a unique ID composes into one
// catalogue alongside the Shipped Packs and participates in normal selection.
func TestExternalPackComposesWithShipped(t *testing.T) {
	t.Parallel()

	externalPack := pack.Definition{ID: "ext", Name: "External"}

	catalog := shipped.ComposeCatalog([]pack.Definition{externalPack})
	packs := catalog.Packs()

	var found bool
	for _, definition := range packs {
		if definition.ID == "ext" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected external pack to compose into the catalogue")
	}

	if _, err := catalog.Registry([]string{"ext"}); err != nil {
		t.Fatalf("expected external pack to select: %v", err)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func mkdirAll(tb testing.TB, path string) {
	tb.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		tb.Fatal(err)
	}
}

func writeFile(tb testing.TB, path string, content string) {
	tb.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		tb.Fatal(err)
	}
}
