package profiles

import (
	"bytes"
	"testing"

	"github.com/BurntSushi/toml"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/profile"
)

/* --------------------------------------- Current Profile -------------------------------------- */

func Current(test *testing.T) (policy profile.Profile) {
	test.Helper()

	policy, err := profile.Load(fixtures.RepoRoot(test))
	if err != nil {
		test.Fatalf("profile.Load: %v", err)
	}

	return policy
}

func RepositoryConfig(test *testing.T) (repository profile.RepositoryConfig) {
	test.Helper()

	return Current(test).Repository
}

func Write(test *testing.T, root string, policy profile.Profile) {
	test.Helper()

	styleGuide := fixtures.ReadFile(test, fixtures.RepoRoot(test), "STYLE.md")
	fixtures.WriteFile(test, root, policy.StyleGuide.Path, styleGuide)
	fixtures.WriteFile(test, root, "style.toml", Render(test, policy))
}

func Render(test *testing.T, policy profile.Profile) (contents string) {
	test.Helper()

	document := struct {
		SchemaVersion int                        `toml:"profile_version"`
		RulePacks     profile.RulePackConfig     `toml:"rule_packs"`
		Repository    profile.RepositoryConfig   `toml:"repository"`
		StyleGuide    profile.StyleGuideConfig   `toml:"styleguide"`
		Imports       profile.ImportsConfig      `toml:"imports"`
		Paths         map[string][]string        `toml:"paths"`
		FileSets      []profile.FileSetConfig    `toml:"file_sets"`
		Language      profile.LanguageConfig     `toml:"language"`
		Naming        profile.NamingConfig       `toml:"naming"`
		ControlPlane  profile.ControlPlaneConfig `toml:"control_plane"`
		Architecture  profile.ArchitectureConfig `toml:"architecture"`
		Rules         []profile.RuleBinding      `toml:"rules"`
	}{
		SchemaVersion: policy.SchemaVersion,
		RulePacks:     policy.RulePacks,
		Repository:    policy.Repository,
		StyleGuide:    policy.StyleGuide,
		Imports:       policy.Imports,
		Paths:         policy.Paths.Classes,
		FileSets:      policy.FileSets,
		Language:      policy.Language,
		Naming:        policy.Naming,
		ControlPlane:  policy.ControlPlane,
		Architecture:  policy.Architecture,
		Rules:         policy.Rules,
	}

	var buffer bytes.Buffer
	encoder := toml.NewEncoder(&buffer)
	if err := encoder.Encode(document); err != nil {
		test.Fatalf("render profile TOML: %v", err)
	}

	return buffer.String()
}
