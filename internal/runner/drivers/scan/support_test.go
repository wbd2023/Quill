package scan

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/policy"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
	"ciphera/tools/internal/rules/vocabulary"
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
		Name: "markdown",
		Include: policy.FileSetInclude{
			Extensions: []string{".md"},
			Files: map[contract.Scope][]string{
				"app": {"STYLE.md"},
			},
			Paths: map[contract.Scope][]string{
				"app":   {"cmd/", "internal/"},
				"tools": {"tools/"},
			},
		},
	})
	goConfig, err := gopolicy.DecodeConfig(config.PackConfigs[builtin.PackGo])
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
			Scope:            contract.Scope("app"),
			WorkingDirectory: ".",
			FormatPaths:      []string{"cmd", "internal"},
			CheckPaths:       []string{"cmd", "internal"},
		},
		{
			Name:             "tools_go",
			Language:         "go",
			Scope:            contract.Scope("tools"),
			WorkingDirectory: "tools",
			FormatPaths:      []string{"cmd", "internal"},
			CheckPaths:       []string{"cmd", "internal"},
		},
	}
	vocabularyConfig, err := vocabulary.DecodeConfig(config.PackConfigs[builtin.PackVocabulary])
	if err != nil {
		t.Fatalf("Decode vocabulary config: %v", err)
	}
	vocabularyConfig.Go.ForbiddenTypeSuffixes = []string{"Repository"}
	vocabularyConfig.Go.PreferredTypeSuffix = "Store"
	vocabularyConfig.Go.ForbiddenIdentifierSuffixes = []string{"Repository"}
	vocabularyConfig.Go.PreferredIdentifierSuffix = "Store"
	config.PackConfigs[builtin.PackVocabulary] = vocabulary.EncodeConfig(vocabularyConfig)
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
	config.PackConfigs[builtin.PackGo] = gopolicy.EncodeConfig(goConfig)

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
