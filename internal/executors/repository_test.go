package executors

import (
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rulepack"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

func TestRunRepositoryScanRuleAcceptsKnownScanner(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.Scope("tools"))

	if _, err := repositoryScanExecutor(
		context,
		repositoryScanSpec(rulepack.ScannerASCII),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanExecutor(ascii): %v", err)
	}
}

func TestRunRepositoryScanRuleRejectsUnknownScanner(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.Scope("all"))

	if _, err := repositoryScanExecutor(
		context,
		repositoryScanSpec("unknown"),
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

	context := testContext(t, fixtureRoot, contract.Scope("all"))
	if result, err := repositoryScanExecutor(
		context,
		repositoryScanSpec(rulepack.ScannerArchitecture),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanExecutor(architecture): %v\n%s", err, result.Output)
	}

	result, err := repositoryScanExecutor(
		context,
		repositoryScanSpec(rulepack.ScannerNaming),
		nil,
	)
	if err == nil {
		t.Fatal("expected alternate naming policy to reject Repository suffixes")
	}

	if !hasDiagnostic(
		result,
		"naming/vocabulary/go-type-suffix",
		"internal/app/services/message_service.go",
		8,
		"use Store not Repository",
	) {
		t.Fatalf("expected alternate naming vocabulary diagnostic, got: %#v", result.Diagnostics)
	}
}

func hasDiagnostic(
	result contract.ExecutionResult,
	code string,
	file string,
	line int,
	messageFragment string,
) (found bool) {
	for _, diagnostic := range result.Diagnostics {
		if diagnostic.Code != code {
			continue
		}
		if file != "" && diagnostic.File != file {
			continue
		}
		if line != 0 && diagnostic.Line != line {
			continue
		}
		if messageFragment != "" && !strings.Contains(diagnostic.Message, messageFragment) {
			continue
		}

		return true
	}

	return false
}

/* -------------------------------------- Fixture Profiles -------------------------------------- */

func alternatePolicyForTest(t *testing.T) (config policy.Config) {
	t.Helper()

	config = profiles.Current(t)
	config.Repository.RootMarkers = []string{"STYLE.md", "style.toml", "ALTROOT"}
	config.Repository.Scopes = map[contract.Scope][]string{
		"app":   {"cmd", "internal"},
		"tools": {"tools"},
		"all":   {"."},
	}
	config.FileSets = replaceFileSet(config.FileSets, policy.FileSetConfig{
		Name:       "markdown",
		Extensions: []string{".md"},
		Files: map[contract.Scope][]string{
			"app": {"STYLE.md"},
		},
		Prefixes: map[contract.Scope][]string{
			"app":   {"cmd/", "internal/"},
			"tools": {"tools/"},
		},
	})
	config.Imports.LocalPrefix = "example.com/altchat"
	config.Paths.Classes = map[string][]string{
		"go_source":        {"cmd/", "internal/"},
		"application_port": {"internal/app/ports/"},
		"concrete_infra":   {"internal/adapters/"},
		"domain":           {"internal/domain/"},
		"domain_errors":    {"internal/domain/errors.go"},
		"test_mocks":       {"internal/testsupport/mocks/"},
	}
	config.Language.Backends = []policy.LanguageBackendConfig{
		{
			Name:        "application_go",
			Language:    "go",
			Scope:       contract.Scope("app"),
			Workdir:     ".",
			FormatPaths: []string{"cmd", "internal"},
			StylePaths:  []string{"cmd", "internal"},
		},
		{
			Name:        "tooling_go",
			Language:    "go",
			Scope:       contract.Scope("tools"),
			Workdir:     "tools",
			FormatPaths: []string{"cmd", "internal"},
			StylePaths:  []string{"tools/cmd", "tools/internal"},
		},
	}
	config.Naming.GoTypeSuffixForbidden = []string{"Repository"}
	config.Naming.GoTypeSuffixPreferred = "Store"
	config.Naming.GoIdentifierSuffixForbidden = []string{"Repository"}
	config.Naming.GoIdentifierSuffixPreferred = "Store"
	config.Naming.GoParameters.ConstructorCategories = replaceConstructorCategory(
		config.Naming.GoParameters.ConstructorCategories,
		policy.GoConstructorCategory{
			Name:        "repository",
			TypeMarkers: []string{"Store"},
		},
	)
	config.Architecture.Layers = []policy.ArchitectureLayer{
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

	return config
}

/* --------------------------------------- Fixture Helpers -------------------------------------- */

func repositoryScanSpec(scanner string) (spec contract.ExecutionSpec) {
	return contract.ExecutionSpec{
		Kind: rulepack.ExecutorRepositoryScan,
		Detail: contract.RepositoryScanExecution{
			Scanner: scanner,
		},
	}
}

func replaceFileSet(
	fileSets []policy.FileSetConfig,
	replacement policy.FileSetConfig,
) (updated []policy.FileSetConfig) {
	updated = append([]policy.FileSetConfig{}, fileSets...)
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
	categories []policy.GoConstructorCategory,
	replacement policy.GoConstructorCategory,
) (updated []policy.GoConstructorCategory) {
	updated = append([]policy.GoConstructorCategory{}, categories...)
	for index, category := range updated {
		if category.Name != replacement.Name {
			continue
		}

		updated[index] = replacement
		return updated
	}

	return append(updated, replacement)
}
