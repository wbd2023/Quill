package executors

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/policy"
)

/* ----------------------------------------- Diagnostics ---------------------------------------- */

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

/* -------------------------------------- Alternate Profile ------------------------------------- */

func alternatePolicyForTest(t *testing.T) (config policy.Config) {
	t.Helper()

	config = profiles.Current(t)
	config.Repository.RootMarkers = []string{"STYLE.md", "style.toml", "ALTROOT"}
	config.Repository.ScopeRoots = map[contract.Scope][]string{
		"app":   {"cmd", "internal"},
		"tools": {"tools"},
		"all":   {"."},
	}
	config.FileSets = replaceFileSet(config.FileSets, policy.FileSetConfig{
		Name:       "markdown",
		Extensions: []string{".md"},
		ExplicitFiles: map[contract.Scope][]string{
			"app": {"STYLE.md"},
		},
		PathPrefixes: map[contract.Scope][]string{
			"app":   {"cmd/", "internal/"},
			"tools": {"tools/"},
		},
	})
	config.Go.LocalImportPrefixes = []string{"example.com/altchat"}
	config.Paths = policy.PathClasses{
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
			WorkDir:     ".",
			FormatPaths: []string{"cmd", "internal"},
			CheckPaths:  []string{"cmd", "internal"},
		},
		{
			Name:        "tooling_go",
			Language:    "go",
			Scope:       contract.Scope("tools"),
			WorkDir:     "tools",
			FormatPaths: []string{"cmd", "internal"},
			CheckPaths:  []string{"cmd", "internal"},
		},
	}
	config.Vocabulary.Go.ForbiddenTypeSuffixes = []string{"Repository"}
	config.Vocabulary.Go.PreferredTypeSuffix = "Store"
	config.Vocabulary.Go.ForbiddenIdentifierSuffixes = []string{"Repository"}
	config.Vocabulary.Go.PreferredIdentifierSuffix = "Store"
	parameters := &config.Go.Parameters
	parameters.ConstructorOrder = replaceParameterGroup(
		parameters.ConstructorOrder,
		policy.GoParameterGroup{
			Name:        "repository",
			TypeMarkers: []string{"Store"},
		},
	)
	config.Go.Architecture.Layers = []policy.GoArchitectureLayer{
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
	groups []policy.GoParameterGroup,
	replacement policy.GoParameterGroup,
) (updated []policy.GoParameterGroup) {
	updated = append([]policy.GoParameterGroup{}, groups...)
	for index, group := range updated {
		if group.Name != replacement.Name {
			continue
		}

		updated[index] = replacement
		return updated
	}

	return append(updated, replacement)
}
