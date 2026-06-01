package validation_test

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/fixture"
	"ciphera/tools/internal/profile/validation"
)

/* --------------------------------------- Profile Version -------------------------------------- */

func TestCheckRequiresCurrentSchemaVersion(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.SchemaVersion = 2
	err := validation.Check(config)
	requireErrorContains(t, err, "version 2")
}

/* ----------------------------------------- Repository ----------------------------------------- */

func TestCheckRejectsUnknownDefaultScope(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.Repository.DefaultScope = "unknown"
	err := validation.Check(config)
	requireErrorContains(t, err, "default_scope")
}

func TestCheckRejectsEmptyRootMarker(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.Repository.RootMarkers = []string{""}
	err := validation.Check(config)
	requireErrorContains(t, err, "root_markers contains an empty marker")
}

func TestCheckRejectsEmptyScopeRoot(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		roots []string
	}{
		{name: "empty root", roots: []string{""}},
		{name: "blank root", roots: []string{"  "}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			config := fixture.Config()

			config.Repository.ScopeRoots[contract.Scope("tools")] = test.roots
			err := validation.Check(config)
			requireErrorContains(
				t,
				err,
				"repository.scope_roots.tools contains an empty root",
			)
		})
	}
}

/* --------------------------------- Path Classes and File Sets --------------------------------- */

func TestCheckAllowsProfileOwnedPathRoles(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.PathRoles["local_policy"] = []string{"internal/local/"}
	if err := validation.Check(config); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestCheckRejectsInvalidPathRole(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.PathRoles["local_policy"] = []string{"internal/local/", " "}
	err := validation.Check(config)
	requireErrorContains(t, err, "path_roles.local_policy")
}

func TestCheckRejectsUnknownFileSetScope(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.FileSets[0].Include.Paths[contract.Scope("unknown")] = []string{"unknown/"}
	err := validation.Check(config)
	requireErrorContains(t, err, "unknown scope")
}

/* --------------------------------------- Packs and Rules -------------------------------------- */

func TestCheckRejectsDuplicateEnabledPacks(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.EnabledPacks = append(config.EnabledPacks, config.EnabledPacks[0])
	err := validation.Check(config)
	requireErrorContains(t, err, "duplicate pack")
}

func TestCheckRejectsDisabledPackConfig(t *testing.T) {
	t.Parallel()

	config := fixture.Config()
	config.PackConfigs = policy.PackConfigs{
		"disabled": policy.PackConfig{"enabled": true},
	}

	err := validation.Check(config)
	requireErrorContains(t, err, "packs.disabled")
}

func TestCheckRejectsEmptyPackConfig(t *testing.T) {
	t.Parallel()

	config := fixture.Config()
	config.PackConfigs = policy.PackConfigs{
		config.EnabledPacks[0]: policy.PackConfig{},
	}

	err := validation.Check(config)
	requireErrorContains(t, err, "must not be empty")
}

func TestCheckRejectsUnknownRuleScope(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.Rules[0].Scope = "unknown"
	err := validation.Check(config)
	requireErrorContains(t, err, "unknown scope")
}

func TestCheckRejectsMalformedRequirementID(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.Rules[0].RequirementIDs = []string{"not-a-requirement-id"}
	err := validation.Check(config)
	requireErrorContains(t, err, "invalid requirement id")
}
