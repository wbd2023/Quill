package coverage

import (
	"testing"

	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/styleguide"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

func loadDocument(t *testing.T) (document styleguide.Document) {
	t.Helper()

	config := profiles.Current(t)
	document, err := styleguide.Load(testutil.RepositoryRoot(t), styleguide.Config{
		Filename: config.StyleGuide.Path,
	})
	if err != nil {
		t.Fatalf("styleguide.Load: %v", err)
	}

	return document
}

func loadPlan(t *testing.T) (plan style.Plan) {
	t.Helper()

	config := profiles.Current(t)
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
		t.Fatalf("profile.Compile: %v", err)
	}

	return compiled.Effective
}

func loadCoverageReport(t *testing.T) (report Report) {
	t.Helper()

	return Build(loadDocument(t), loadPlan(t).Rules)
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
