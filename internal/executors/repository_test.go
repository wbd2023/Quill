package executors

import (
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

func TestRunRepositoryScanRuleAcceptsKnownScanner(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.ScopeTools)

	if _, err := repositoryScanExecutor(
		context,
		contract.ExecutionSpec{Scanner: rulepack.RepositoryScannerASCII},
		nil,
	); err != nil {
		t.Fatalf("repositoryScanExecutor(ascii): %v", err)
	}
}

func TestRunRepositoryScanRuleRejectsUnknownScanner(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.ScopeAll)

	if _, err := repositoryScanExecutor(
		context,
		contract.ExecutionSpec{Scanner: "unknown"},
		nil,
	); err == nil {
		t.Fatal("expected unknown scanner to be rejected")
	}
}

func TestRunRepositoryScanRuleSupportsAlternateProfile(t *testing.T) {
	fixtureRoot := t.TempDir()
	alternateProfile := alternatePolicyForTest(t)
	profiles.Write(t, fixtureRoot, alternateProfile)
	fixtures.WriteFile(t, fixtureRoot, "ALTROOT", "")
	fixtures.WriteFile(
		t,
		fixtureRoot,
		"go.mod",
		"module example.com/altchat\n\ngo 1.24.5\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "errors.go"),
		"package domain\n\nvar ErrMissing = error(nil)\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "app", "ports", "message_store.go"),
		"package ports\n\ntype Message"+"Store interface { ListMessages() }\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "app", "services", "message_service.go"),
		"package services\n\n"+
			"import (\n"+
			"\t\"example.com/altchat/internal/app/ports\"\n"+
			"\t\"example.com/altchat/internal/domain\"\n"+
			")\n\n"+
			"type Message"+"Repository interface {\n"+
			"\tListMessages() []domain.Message\n"+
			"}\n\n"+
			"type MessageService struct {\n"+
			"\tstore ports.Message"+"Store\n"+
			"}\n",
	)
	fixtures.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "message.go"),
		"package domain\n\ntype Message struct{}\n",
	)

	context := testContext(t, fixtureRoot, contract.ScopeAll)
	if output, err := repositoryScanExecutor(
		context,
		contract.ExecutionSpec{Scanner: rulepack.RepositoryScannerArchitecture},
		nil,
	); err != nil {
		t.Fatalf("repositoryScanExecutor(architecture): %v\n%s", err, output)
	}

	output, err := repositoryScanExecutor(
		context,
		contract.ExecutionSpec{Scanner: rulepack.RepositoryScannerNaming},
		nil,
	)
	if err == nil {
		t.Fatal("expected alternate naming policy to reject Repository suffixes")
	}

	if !strings.Contains(output, "use Store not Repository") {
		t.Fatalf("expected alternate naming vocabulary in output, got:\n%s", output)
	}
}

/* -------------------------------------- Fixture Profiles -------------------------------------- */

func alternatePolicyForTest(t *testing.T) (policy profile.Profile) {
	t.Helper()

	policy = profiles.Current(t)
	policy.Repository.RootMarkers = []string{"STYLE.md", "style.toml", "ALTROOT"}
	policy.Repository.AppScanRoots = []string{"cmd", "internal"}
	policy.Repository.ToolsScanRoots = []string{"tools"}
	policy.FileSets = replaceFileSet(policy.FileSets, profile.FileSetConfig{
		Name:          "markdown",
		Extensions:    []string{".md"},
		AppFiles:      []string{"STYLE.md"},
		AppPrefixes:   []string{"cmd/", "internal/"},
		ToolsPrefixes: []string{"tools/"},
	})
	policy.Imports.LocalPrefix = "example.com/altchat"
	policy.Paths.Classes = map[string][]string{
		rulepack.PathClassApp:             {"cmd/", "internal/"},
		rulepack.PathClassApplicationPort: {"internal/app/ports/"},
		rulepack.PathClassConcreteInfra:   {"internal/adapters/"},
		rulepack.PathClassDomain:          {"internal/domain/"},
		rulepack.PathClassDomainErrors:    {"internal/domain/errors.go"},
		rulepack.PathClassTestMocks:       {"internal/testsupport/mocks/"},
	}
	policy.Language.Backends = []profile.LanguageBackendConfig{
		{
			Name:        "go_app",
			Language:    "go",
			Workdir:     ".",
			FormatPaths: []string{"cmd", "internal"},
			StylePaths:  []string{"cmd", "internal"},
		},
		{
			Name:        "go_tools",
			Language:    "go",
			Workdir:     "tools",
			FormatPaths: []string{"cmd", "internal"},
			StylePaths:  []string{"tools/cmd", "tools/internal"},
		},
	}
	policy.Naming.GoTypeSuffixForbidden = []string{"Repository"}
	policy.Naming.GoTypeSuffixPreferred = "Store"
	policy.Naming.GoIdentifierSuffixForbidden = []string{"Repository"}
	policy.Naming.GoIdentifierSuffixPreferred = "Store"
	policy.Naming.GoParameters.ConstructorCategories = replaceConstructorCategory(
		policy.Naming.GoParameters.ConstructorCategories,
		profile.GoConstructorCategory{
			Name:        "repository",
			TypeMarkers: []string{"Store"},
		},
	)
	policy.Architecture.Layers = []profile.ArchitectureLayer{
		{
			Name:         "domain",
			PackageRoots: []string{"internal/domain"},
			MayImport:    []string{"domain"},
		},
		{
			Name:         "port",
			PackageRoots: []string{"internal/app/ports"},
			MayImport:    []string{"domain", "port"},
		},
		{
			Name:         "service",
			PackageRoots: []string{"internal/app/services"},
			MayImport:    []string{"domain", "port", "service"},
		},
		{
			Name:         "adapter",
			PackageRoots: []string{"internal/adapters"},
			MayImport:    []string{"domain", "port", "service", "adapter"},
		},
		{
			Name:         "cmd",
			PackageRoots: []string{"cmd"},
			MayImport:    []string{"service", "adapter"},
		},
	}

	return policy
}

/* --------------------------------------- Fixture Helpers -------------------------------------- */

func replaceFileSet(
	fileSets []profile.FileSetConfig,
	replacement profile.FileSetConfig,
) (updated []profile.FileSetConfig) {
	updated = append([]profile.FileSetConfig{}, fileSets...)
	for index, fileSet := range updated {
		if fileSet.Name != replacement.Name {
			continue
		}

		updated[index] = replacement
		return updated
	}

	return append(updated, replacement)
}

func replaceConstructorCategory(
	categories []profile.GoConstructorCategory,
	replacement profile.GoConstructorCategory,
) (updated []profile.GoConstructorCategory) {
	updated = append([]profile.GoConstructorCategory{}, categories...)
	for index, category := range updated {
		if category.Name != replacement.Name {
			continue
		}

		updated[index] = replacement
		return updated
	}

	return append(updated, replacement)
}
