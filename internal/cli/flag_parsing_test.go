package cli

import (
	"strings"
	"testing"

	"ciphera/tools/internal/testutil"
)

/* ------------------------------------------- Parsing ------------------------------------------ */

func TestParseCheckOptionsRejectsInvalidMode(t *testing.T) {
	_, err := parseCheckOptions([]string{"--mode", "invalid"})
	if err == nil {
		t.Fatal("expected invalid mode error")
	}
}

func TestParseCheckOptionsAcceptsConfiguredScopeName(t *testing.T) {
	options, err := parseCheckOptions([]string{"--scope", "custom"})
	if err != nil {
		t.Fatalf("parseCheckOptions: %v", err)
	}

	if options.scope != "custom" {
		t.Fatalf("scope = %q, want custom", options.scope)
	}
}

func TestParseCheckOptionsRejectsPositionalArguments(t *testing.T) {
	_, err := parseCheckOptions([]string{"extra"})
	if err == nil {
		t.Fatal("expected positional argument error")
	}

	if !strings.Contains(err.Error(), "extra") {
		t.Fatalf("expected positional argument in error, got %q", err)
	}
}

func TestParseDoctorOptionsAcceptsRepoRoot(t *testing.T) {
	options, err := parseDoctorOptions([]string{"--repo-root", "."})
	if err != nil {
		t.Fatalf("expected repo-root to parse, got %v", err)
	}

	if options.repoRoot == "" {
		t.Fatal("expected absolute repo-root")
	}
}

func TestParseInstallOptionsRejectsPositionalArguments(t *testing.T) {
	_, err := parseInstallOptions([]string{"extra"})
	if err == nil {
		t.Fatal("expected positional argument error")
	}
}

func TestParseFixOptionsAcceptsScope(t *testing.T) {
	options, err := parseFixOptions([]string{"--scope", "tools"})
	if err != nil {
		t.Fatalf("parseFixOptions: %v", err)
	}

	if options.scope != "tools" {
		t.Fatalf("expected tools scope, got %q", options.scope)
	}
}

func TestParseCoverageOptionsAcceptsVerbose(t *testing.T) {
	options, err := parseCoverageOptions([]string{"--verbose"})
	if err != nil {
		t.Fatalf("parseCoverageOptions: %v", err)
	}

	if !options.verbose {
		t.Fatal("expected verbose coverage output to be enabled")
	}
}

func TestParseCheckOptionsAcceptsJSONFormat(t *testing.T) {
	options, err := parseCheckOptions([]string{"--format", "json"})
	if err != nil {
		t.Fatalf("parseCheckOptions: %v", err)
	}

	if options.format != "json" {
		t.Fatalf("expected json format, got %q", options.format)
	}
}

func TestParseCoverageOptionsRejectsInvalidFormat(t *testing.T) {
	_, err := parseCoverageOptions([]string{"--format", "yaml"})
	if err == nil {
		t.Fatal("expected invalid format error")
	}
}

func TestResolveRepoRootAutoDetectsRepository(t *testing.T) {
	repositoryRoot, err := resolveRepoRoot("")
	if err != nil {
		t.Fatalf("resolveRepoRoot: %v", err)
	}

	if repositoryRoot != testutil.RepositoryRoot(t) {
		t.Fatalf("unexpected repo root %q", repositoryRoot)
	}
}

func TestFindRepoRootRejectsMissingRepository(t *testing.T) {
	missingRoot := t.TempDir()

	_, err := findRepoRoot(missingRoot)
	if err == nil {
		t.Fatal("expected missing repository error")
	}
}

func TestFindRepoRootRejectsLegacyRootWithoutStyleProfile(t *testing.T) {
	legacyRoot := t.TempDir()
	testutil.WriteFile(t, legacyRoot, "STYLE.md", "# Style\n")

	testutil.WriteFile(t, legacyRoot, "tools/go.mod", "module example.com/tools\n")

	if _, err := findRepoRoot(legacyRoot); err == nil {
		t.Fatal("expected legacy repo root without style.toml to be rejected")
	}
}
