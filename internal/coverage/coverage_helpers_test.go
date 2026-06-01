package coverage

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/profile/effective"
	"ciphera/tools/internal/styleguide"
)

func loadDocument(t *testing.T) (document styleguide.Document) {
	t.Helper()

	config := profiles.Current(t)
	document, err := styleguide.Load(fixtures.RepositoryRoot(t), styleguide.Config{
		Filename: config.StyleGuide.Path,
		IDScheme: config.StyleGuide.IDScheme,
	})
	if err != nil {
		t.Fatalf("styleguide.Load: %v", err)
	}

	return document
}

func loadEffectiveConfig(t *testing.T) (compiled contract.EffectiveConfig) {
	t.Helper()

	config := profiles.Current(t)
	registry, err := builtin.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	config, err = effective.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("effective.ResolvePacks: %v", err)
	}

	compiled, err = profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("profile.Compile: %v", err)
	}

	return compiled
}

func loadCoverageReport(t *testing.T) (report Report) {
	t.Helper()

	return Build(loadDocument(t), loadEffectiveConfig(t).Rules)
}

func coverageRequirementByID(
	report Report,
	requirementID string,
) (requirement Requirement, found bool) {
	for _, requirement := range report.Requirements {
		if requirement.ID == requirementID {
			return requirement, true
		}
	}

	return Requirement{}, false
}
