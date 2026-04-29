package coverage

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/requirementid"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/styleguide"
)

func loadDocument(t *testing.T) (document styleguide.Document) {
	t.Helper()

	config := profiles.Current(t)
	document, err := styleguide.Load(fixtures.RepoRoot(t), styleguide.Config{
		Filename:            config.StyleGuide.Path,
		RequirementIDScheme: requirementid.Scheme(config.StyleGuide.RequirementIDScheme),
	})
	if err != nil {
		t.Fatalf("styleguide.Load: %v", err)
	}

	return document
}

func loadEffectiveConfig(t *testing.T) (effective contract.EffectiveConfig) {
	t.Helper()

	config := profiles.Current(t)
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err = profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("profile.Compile: %v", err)
	}

	return effective
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
