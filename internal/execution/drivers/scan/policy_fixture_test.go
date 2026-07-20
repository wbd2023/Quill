package scan

import (
	"slices"
	"testing"

	"github.com/wbd2023/Quill/internal/checks/gopolicy"
	"github.com/wbd2023/Quill/internal/checks/vocabularypolicy"
	"github.com/wbd2023/Quill/internal/pack/shipped/golang"
	"github.com/wbd2023/Quill/internal/pack/shipped/vocabulary"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

/* --------------------------------------- Policy Fixture --------------------------------------- */

func buildScanDriverPolicyFixture(t *testing.T) (config policy.Config) {
	t.Helper()

	config = profiles.Current(t)
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
