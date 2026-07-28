package bindings

import (
	"context"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	goenv "github.com/wbd2023/quill/internal/ecosystem/golang"
	"github.com/wbd2023/quill/internal/ecosystem/node"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
	"github.com/wbd2023/quill/internal/pack/shipped/text"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/pack/shipped/vocabulary"
	vocabularypolicy "github.com/wbd2023/quill/internal/pack/shipped/vocabulary/policy"
	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
	"github.com/wbd2023/quill/internal/workspace"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

// These integration tests exercise the real Pack-owned scanner closures through the generic
// repository-scan dispatch. They relocated here from the dismantled execution/drivers/scan package
// because the scanners now live in the Pack bindings children and only the composition root can
// reach both the bindings and the generic driver set.

func TestRepositoryScanDriverAcceptsKnownScanner(t *testing.T) {
	runCtx := testContext(t, testutil.RepositoryRoot(t), style.Scope("all"))

	if _, err := scanDriver()(
		context.Background(),
		runCtx,
		repositoryScanSpec(text.PackID, text.ScannerASCII),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanDriver(ascii): %v", err)
	}
}

func TestRepositoryScanDriverRejectsUnknownScanner(t *testing.T) {
	runCtx := testContext(t, testutil.RepositoryRoot(t), style.Scope("all"))

	_, err := scanDriver()(
		context.Background(),
		runCtx,
		repositoryScanSpec("", "unknown"),
		nil,
	)
	if err == nil {
		t.Fatal("expected unknown scanner to be rejected")
	}

	if !strings.Contains(err.Error(), `"unknown"`) {
		t.Fatalf("error = %q, want scanner ID", err)
	}
}

func TestRepositoryScanDriverSupportsAlternateProfile(t *testing.T) {
	fixtureRoot := t.TempDir()
	alternateProfile := buildAlternateProfile(t)
	profiles.Write(t, fixtureRoot, alternateProfile)
	testutil.WriteFile(t, fixtureRoot, "ALTROOT", "")
	testutil.WriteFile(
		t,
		fixtureRoot,
		"go.mod",
		"module example.com/altchat\n\ngo 1.24.5\n",
	)
	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "errors.go"),
		"package domain\n\nvar ErrMissing = error(nil)\n",
	)
	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "app", "ports", "message_store.go"),
		"package ports\n\ntype Message"+"Store interface { ListMessages() }\n",
	)
	testutil.WriteFile(
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
	testutil.WriteFile(
		t,
		fixtureRoot,
		filepath.Join("internal", "domain", "message.go"),
		"package domain\n\ntype Message struct{}\n",
	)

	runCtx := testContext(t, fixtureRoot, style.Scope("all"))
	if _, err := scanDriver()(
		context.Background(),
		runCtx,
		repositoryScanSpec(golang.PackID, golang.ScannerArchitecture),
		nil,
	); err != nil {
		t.Fatalf("repositoryScanDriver(architecture): %v", err)
	}
	result, err := scanDriver()(
		context.Background(),
		runCtx,
		repositoryScanSpec(vocabulary.PackID, vocabulary.ScannerVocabulary),
		nil,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
		t.Fatal("expected alternate vocabulary policy to reject Repository suffixes")
	}

	if !hasDiagnosticMatching(
		result,
		"vocabulary/project-terms/go-type-suffix",
		"internal/app/services/message_service.go",
		8,
		"must be Store",
	) {
		t.Fatalf("expected alternate vocabulary diagnostic, got: %#v", result.Diagnostics)
	}
}

/* --------------------------------------- Target Commands -------------------------------------- */

// TestGolangciTargetCommandPassesCleanRepository relocated from execution/drivers/target: it drives
// the real Go golangci target command through the generic target-command dispatch over a clean
// repository with stub goimports and golangci-lint executables.
func TestGolangciTargetCommandPassesCleanRepository(t *testing.T) {
	repoRoot := t.TempDir()
	profiles.Write(t, repoRoot, profiles.Self(t))
	testutil.WriteFile(t, repoRoot, "cmd/quill/main.go", "package main\n\nfunc main() {}\n")
	testutil.WriteFile(t, repoRoot, "internal/example/example.go", "package example\n")
	writeExecutable(t, repoRoot, "goimports")
	writeExecutable(t, repoRoot, "golangci-lint")

	runCtx := testContext(t, repoRoot, style.Scope("all"))

	job := style.TargetCommandJob{
		PackID:   golang.PackID,
		ToolIDs:  []string{tool.Go, tool.Goimports, tool.GolangciLint},
		Action:   golang.TargetActionGolangci,
		Language: golang.Language,
		Targets:  []string{"go"},
	}

	result, err := targetCommandDriver()(context.Background(), runCtx, job, nil)
	if err != nil {
		t.Fatalf("golangciDriver(all): %v", err)
	}

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected repository lint diagnostics: %#v", result.Diagnostics)
	}
}

func writeExecutable(t *testing.T, repoRoot string, name string) {
	t.Helper()

	path := testutil.WriteFile(
		t,
		repoRoot,
		filepath.Join(".cache", "quill", "bin", name),
		"#!/bin/sh\nexit 0\n",
	)
	if err := os.Chmod(path, 0o755); err != nil {
		t.Fatalf("make %s executable: %v", name, err)
	}
}

/* ---------------------------------------- Test Helpers ---------------------------------------- */

func repositoryScanSpec(packID string, scanner string) (job style.Job) {
	return style.RepositoryScanExecution{
		PackID:  packID,
		Scanner: scanner,
	}
}

func scanDriver() (driver execution.Driver) {
	return drivers.CheckDrivers(Build()).RepositoryScan
}

func targetCommandDriver() (driver execution.Driver) {
	return drivers.CheckDrivers(Build()).TargetCommand
}

func testContext(
	t *testing.T,
	repoRoot string,
	scope style.Scope,
) (context execution.RunContext) {
	t.Helper()

	config, err := profile.Load(repoRoot)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := shipped.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	config, err = pack.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}

	compiled, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	layout := workspace.NewLayout(repoRoot)
	path := layout.BuildPath(node.BinaryDirectory(layout))
	toolEnvironment := map[string]string{"PATH": path}
	goEnvironment := goenv.Environment(layout, path)
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDirectory(), "golangci")

	return execution.NewRunContext(
		repoRoot,
		scope,
		compiled.Profile,
		compiled.Effective,
		registry.ToolCapabilities(),
		toolEnvironment,
		goEnvironment,
	)
}

func hasDiagnosticMatching(
	result style.ExecutionResult,
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
		if line != 0 && diagnostic.Range.Start.Line != line {
			continue
		}
		if messageFragment != "" && !strings.Contains(diagnostic.Message, messageFragment) {
			continue
		}

		return true
	}

	return false
}

/* --------------------------------------- Policy Fixture --------------------------------------- */

func buildAlternateProfile(t *testing.T) (config policy.Profile) {
	t.Helper()

	config = profiles.Self(t)
	if !slices.Contains(config.EnabledPacks, vocabulary.PackID) {
		config.EnabledPacks = append(config.EnabledPacks, vocabulary.PackID)
	}
	config.Repository.RootMarkers = []string{"STYLE.md", "quill.toml", "ALTROOT"}
	config.Repository.ScopeRoots = map[style.Scope][]string{
		"app":   {"cmd", "internal"},
		"tools": {"tools"},
		"all":   {"."},
	}
	config.FileSets = replaceFileSet(config.FileSets, policy.FileSetConfig{
		Name: "markdown",
		Include: policy.FileSetInclude{
			Extensions: []string{".md"},
			Files: map[style.Scope][]string{
				"app": {"STYLE.md"},
			},
			Paths: map[style.Scope][]string{
				"app":   {"cmd/", "internal/"},
				"tools": {"tools/"},
			},
		},
	})
	goConfig, err := gopolicy.DecodeConfig(config.PackConfigs[golang.PackID])
	if err != nil {
		t.Fatalf("Decode Go config: %v", err)
	}
	goConfig.LocalImportPrefixes = []string{"example.com/altchat"}
	config.PathRoles = policy.PathRoles{
		"go_source":        {"cmd/", "internal/"},
		"application_port": {"internal/app/ports/"},
		"concrete_infra":   {"internal/adapters/"},
		"domain":           {"internal/domain/"},
		"domain_errors":    {"internal/domain/errors.go"},
		"test_mocks":       {"internal/testsupport/mocks/"},
	}
	config.Targets = []policy.TargetConfig{
		{
			Name:             "app_go",
			Language:         "go",
			Scope:            style.Scope("app"),
			WorkingDirectory: ".",
			FormatPaths:      []string{"cmd", "internal"},
			CheckPaths:       []string{"cmd", "internal"},
		},
		{
			Name:             "tools_go",
			Language:         "go",
			Scope:            style.Scope("tools"),
			WorkingDirectory: "tools",
			FormatPaths:      []string{"cmd", "internal"},
			CheckPaths:       []string{"cmd", "internal"},
		},
	}
	vocabularyConfig := vocabularypolicy.Config{
		Go: vocabularypolicy.GoConfig{
			TypeSuffixes:       map[string][]string{"Store": {"Repository"}},
			IdentifierSuffixes: map[string][]string{"Store": {"Repository"}},
		},
	}
	config.PackConfigs[vocabulary.PackID] = vocabularypolicy.EncodeConfig(vocabularyConfig)
	parameters := &goConfig.Constructors
	parameters.ParameterOrder = replaceParameterGroup(
		parameters.ParameterOrder,
		gopolicy.ParameterGroup{
			Name:             "repository",
			TypeNameSuffixes: []string{"Store"},
		},
	)
	goConfig.Architecture.Layers = []gopolicy.ArchitectureLayer{
		{
			Name:          "domain",
			PackageRoots:  []string{"internal/domain"},
			AllowedLayers: []string{"domain"},
		},
		{
			Name:          "port",
			PackageRoots:  []string{"internal/app/ports"},
			AllowedLayers: []string{"domain", "port"},
		},
		{
			Name:          "service",
			PackageRoots:  []string{"internal/app/services"},
			AllowedLayers: []string{"domain", "port", "service"},
		},
		{
			Name:          "adapter",
			PackageRoots:  []string{"internal/adapters"},
			AllowedLayers: []string{"domain", "port", "service", "adapter"},
		},
		{
			Name:          "cmd",
			PackageRoots:  []string{"cmd"},
			AllowedLayers: []string{"service", "adapter"},
		},
	}
	config.PackConfigs[golang.PackID] = gopolicy.EncodeConfig(goConfig)

	return config
}

/* --------------------------------------- Config Updates --------------------------------------- */

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

func replaceParameterGroup(
	groups []gopolicy.ParameterGroup,
	replacement gopolicy.ParameterGroup,
) (updated []gopolicy.ParameterGroup) {
	updated = append([]gopolicy.ParameterGroup{}, groups...)
	for index, group := range updated {
		if group.Name != replacement.Name {
			continue
		}

		updated[index] = replacement
		return updated
	}

	return append(updated, replacement)
}
