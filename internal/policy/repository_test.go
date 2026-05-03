package policy_test

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

const (
	scopeAll         contract.Scope = "all"
	scopeApp         contract.Scope = "app"
	scopeCommand     contract.Scope = "command"
	scopeCommandLine contract.Scope = "command_line"
	scopeNested      contract.Scope = "nested"
	scopeTools       contract.Scope = "tools"
)

/* ------------------------------------------- Scopes ------------------------------------------- */

func TestRepositoryHasScope(t *testing.T) {
	repository := testRepository()

	if !repository.HasScope(scopeApp) {
		t.Fatalf("expected repository to contain scope %q", scopeApp)
	}
	if repository.HasScope("missing") {
		t.Fatalf("expected repository not to contain missing scope")
	}
}

func TestRepositoryResolveScopeRoots(t *testing.T) {
	repositoryRoot := filepath.Join("workspace", "repo")
	repository := testRepository()

	roots := repository.ResolveScopeRoots(repositoryRoot, scopeApp)
	requireEqual(t, []string{
		filepath.Join(repositoryRoot, "cmd"),
		filepath.Join(repositoryRoot, "internal"),
	}, roots)

	roots = repository.ResolveScopeRoots(repositoryRoot, scopeAll)
	requireEqual(t, []string{repositoryRoot}, roots)
}

func TestRepositoryResolveScopeRootsNormalisesConfiguredRoots(t *testing.T) {
	repositoryRoot := filepath.Join("workspace", "repo")
	repository := testRepository()
	normalisedScope := contract.Scope("normalised")

	repository.ScopeRoots[normalisedScope] = []string{" ./cmd/../tools/ "}
	roots := repository.ResolveScopeRoots(repositoryRoot, normalisedScope)
	requireEqual(t, []string{filepath.Join(repositoryRoot, "tools")}, roots)
}

func TestRepositoryHasScopeOverlap(t *testing.T) {
	repository := testRepository()

	cases := []struct {
		name     string
		scope    contract.Scope
		other    contract.Scope
		expected bool
	}{
		{name: "global scope", scope: scopeAll, other: scopeTools, expected: true},
		{name: "same root", scope: scopeApp, other: scopeCommand, expected: true},
		{name: "nested root", scope: scopeCommand, other: scopeNested, expected: true},
		{name: "sibling prefix", scope: scopeCommand, other: scopeCommandLine, expected: false},
		{name: "separate roots", scope: scopeApp, other: scopeTools, expected: false},
		{name: "missing scope", scope: "missing", other: scopeTools, expected: false},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			overlap := repository.HasScopeOverlap(test.scope, test.other)
			if overlap != test.expected {
				t.Fatalf("unexpected scope overlap %t", overlap)
			}
		})
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func testRepository() (repository policy.RepositoryConfig) {
	return policy.RepositoryConfig{
		ScopeRoots: map[contract.Scope][]string{
			scopeAll:         {"."},
			scopeApp:         {"cmd", "internal"},
			scopeCommand:     {"cmd"},
			scopeCommandLine: {"cmdline"},
			scopeNested:      {"cmd/client"},
			scopeTools:       {"tools"},
		},
	}
}
