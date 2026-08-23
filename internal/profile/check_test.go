package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
	"github.com/wbd2023/quill/internal/style"
)

/* --------------------------------------- Profile Version -------------------------------------- */

func TestCheckRequiresCurrentSchemaVersion(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.SchemaVersion = 2
	err := profile.Validate(config)
	requireErrorContains(t, err, "version 2")
}

/* ----------------------------------------- Repository ----------------------------------------- */

func TestCheckRejectsUnknownDefaultScope(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.Repository.DefaultScope = "unknown"
	err := profile.Validate(config)
	requireErrorContains(t, err, "default_scope")
}

func TestCheckRejectsEmptyRootMarker(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.Repository.RootMarkers = []string{""}
	err := profile.Validate(config)
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

			config := profiletest.Config()

			config.Repository.ScopeRoots[style.Scope("tools")] = test.roots
			err := profile.Validate(config)
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

	config := profiletest.Config()

	config.PathRoles["local_policy"] = []string{"internal/local/"}
	if err := profile.Validate(config); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestCheckRejectsInvalidPathRole(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.PathRoles["local_policy"] = []string{"internal/local/", " "}
	err := profile.Validate(config)
	requireErrorContains(t, err, "path_roles.local_policy")
}

func TestCheckRejectsUnknownFileSetScope(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.FileSets[0].Include.Paths[style.Scope("unknown")] = []string{"unknown/"}
	err := profile.Validate(config)
	requireErrorContains(t, err, "unknown scope")
}

/* --------------------------------------- Packs and Rules -------------------------------------- */

func TestCheckRejectsDuplicateEnabledPacks(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.EnabledPacks = append(config.EnabledPacks, config.EnabledPacks[0])
	err := profile.Validate(config)
	requireErrorContains(t, err, "duplicate pack")
}

func TestCheckRejectsReservedEnabledPack(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.EnabledPacks = []string{profile.EnabledPacksKey}

	requireErrorContains(t, profile.Validate(config), "reserved pack")
	_, err := profile.Format(config)
	requireErrorContains(t, err, "reserved pack")
}

func TestCheckRejectsDisabledPackPolicy(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.PackPolicies = profile.PackPolicies{
		"disabled": profile.PackPolicy{"enabled": true},
	}

	err := profile.Validate(config)
	requireErrorContains(t, err, "packs.disabled")
}

func TestCheckAllowsEmptyPackPolicy(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.PackPolicies = profile.PackPolicies{
		config.EnabledPacks[0]: {},
	}

	if err := profile.Validate(config); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestCheckAllowsDistinctRulesToBindSameRequirement(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	overlappingBinding := config.Rules[0]
	overlappingBinding.RuleID = "test/overlapping-rule"
	config.Rules = append(config.Rules, overlappingBinding)

	if err := profile.Validate(config); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestCheckRejectsUnknownRuleScope(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.Rules[0].Scope = "unknown"
	err := profile.Validate(config)
	requireErrorContains(t, err, "unknown scope")
}

func TestCheckRejectsMalformedRequirementID(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.Rules[0].RequirementIDs = []string{"not-a-requirement-id"}
	err := profile.Validate(config)
	requireErrorContains(t, err, "invalid requirement id")
}
