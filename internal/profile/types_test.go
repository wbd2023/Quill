package profile_test

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"

	"github.com/google/go-cmp/cmp"
)

const (
	scopeAll         style.Scope = "all"
	scopeApp         style.Scope = "app"
	scopeCommand     style.Scope = "command"
	scopeCommandLine style.Scope = "command_line"
	scopeNested      style.Scope = "nested"
	scopeTools       style.Scope = "tools"
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
	root := filepath.Join("workspace", "repo")
	repository := testRepository()

	roots := repository.ResolveScopeRoots(root, scopeApp)
	requireEqual(t, []string{
		filepath.Join(root, "cmd"),
		filepath.Join(root, "internal"),
	}, roots)

	roots = repository.ResolveScopeRoots(root, scopeAll)
	requireEqual(t, []string{root}, roots)
}

func TestRepositoryResolveScopeRootsNormalisesConfiguredRoots(t *testing.T) {
	root := filepath.Join("workspace", "repo")
	repository := testRepository()
	normalisedScope := style.Scope("normalised")

	repository.ScopeRoots[normalisedScope] = []string{" ./cmd/../tools/ "}
	roots := repository.ResolveScopeRoots(root, normalisedScope)
	requireEqual(t, []string{filepath.Join(root, "tools")}, roots)
}

func TestRepositoryHasScopeOverlap(t *testing.T) {
	repository := testRepository()

	cases := []struct {
		name     string
		scope    style.Scope
		other    style.Scope
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

/* ------------------------------------------ File Sets ----------------------------------------- */

func TestFileSetsLookup(t *testing.T) {
	fileSets := profile.FileSets{
		{Name: "markdown", Include: profile.FileSetInclude{Extensions: []string{".md"}}},
	}

	fileSet, found := fileSets.Lookup("markdown")
	if !found {
		t.Fatalf("expected file set lookup to find markdown")
	}

	requireEqual(t, profile.FileSetConfig{
		Name: "markdown",
		Include: profile.FileSetInclude{
			Extensions: []string{".md"},
		},
	}, fileSet)

	_, found = fileSets.Lookup("missing")
	if found {
		t.Fatalf("expected missing file set lookup to fail")
	}
}

/* ------------------------------------------- Targets ------------------------------------------ */

func TestTargetConfigsLookup(t *testing.T) {
	var targets profile.TargetConfigs
	targets = append(targets, profile.TargetConfig{
		Name:     "tools_go",
		Language: "go",
		Scope:    "tools",
	})

	target, found := targets.Lookup("tools_go")
	if !found {
		t.Fatalf("expected target lookup to find tools_go")
	}

	requireEqual(t, profile.TargetConfig{
		Name:     "tools_go",
		Language: "go",
		Scope:    style.Scope("tools"),
	}, target)

	_, found = targets.Lookup("missing")
	if found {
		t.Fatalf("expected missing target lookup to fail")
	}
}

/* -------------------------------------------- Tools ------------------------------------------- */

func TestPinnedToolsLookup(t *testing.T) {
	tools := profile.PinnedTools{
		{ID: "go", Version: "1.24.5", TimeoutSeconds: 30},
	}

	pinnedTool, found := tools.Lookup("go")
	if !found {
		t.Fatalf("expected pinned tool lookup to find go")
	}

	requireEqual(t, profile.PinnedTool{
		ID:             "go",
		Version:        "1.24.5",
		TimeoutSeconds: 30,
	}, pinnedTool)

	_, found = tools.Lookup("missing")
	if found {
		t.Fatalf("expected missing pinned tool lookup to fail")
	}
}

/* ----------------------------------------- Path Roles ----------------------------------------- */

func TestPathRolesLookupPatterns(t *testing.T) {
	roles := profile.PathRoles{
		"go_source": {"cmd/", "internal/"},
	}

	patterns := roles.LookupPatterns("go_source")
	requireEqual(t, []string{"cmd/", "internal/"}, patterns)

	patterns[0] = "mutated/"
	requireEqual(t, []string{"cmd/", "internal/"}, roles.LookupPatterns("go_source"))
}

func TestPathRolesLookupPatternsHandlesMissingClasses(t *testing.T) {
	cases := []struct {
		name  string
		roles profile.PathRoles
	}{
		{name: "nil roles", roles: nil},
		{
			name:  "unknown role",
			roles: profile.PathRoles{"markdown": {".md"}},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			requireEqual(t, []string(nil), test.roles.LookupPatterns("go_source"))
		})
	}
}

/* ---------------------------------------- Pack Policies --------------------------------------- */

func TestPackPolicyCloneCopiesArraysOfTables(t *testing.T) {
	t.Parallel()

	config := profile.PackPolicy{
		"tables": []map[string]any{
			{"name": "first", "values": []string{"a"}},
		},
	}

	clone := config.Clone()
	config["tables"].([]map[string]any)[0]["name"] = "changed"
	config["tables"].([]map[string]any)[0]["values"].([]string)[0] = "b"

	cloneTable := clone["tables"].([]map[string]any)[0]
	requireEqual(t, "first", cloneTable["name"])
	requireEqual(t, "a", cloneTable["values"].([]string)[0])
}

func requireEqual[T any](t *testing.T, want T, got T) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("unexpected value (-expected +actual):\n%s", diff)
	}
}

func testRepository() (repository profile.RepositoryConfig) {
	return profile.RepositoryConfig{
		ScopeRoots: map[style.Scope][]string{
			scopeAll:         {"."},
			scopeApp:         {"cmd", "internal"},
			scopeCommand:     {"cmd"},
			scopeCommandLine: {"cmdline"},
			scopeNested:      {"cmd/client"},
			scopeTools:       {"tools"},
		},
	}
}
