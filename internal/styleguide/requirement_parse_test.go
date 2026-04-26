package styleguide

import "testing"

func TestParseStyleRequirementID(t *testing.T) {
	requirementID, _, found := parseRequirementText(
		"* `[3.2.ctx-first]` `ctx context.Context` MUST be the first parameter when present.",
		"",
		RequirementIDFormatSectionSlug,
	)
	if !found {
		t.Fatal("expected requirement ID to parse")
	}

	if requirementID != "3.2.ctx-first" {
		t.Fatalf("unexpected requirement ID %q", requirementID)
	}
}

func TestParseStyleRequirement(t *testing.T) {
	requirementID, requirementText, found := parseRequirementText(
		"* `[3.2.ctx-first]` `ctx context.Context` MUST be the first parameter when present.",
		"",
		RequirementIDFormatSectionSlug,
	)
	if !found {
		t.Fatal("expected requirement to parse")
	}

	if requirementID != "3.2.ctx-first" {
		t.Fatalf("unexpected requirement ID %q", requirementID)
	}

	if requirementText == "" {
		t.Fatal("expected requirement text")
	}
}

func TestRequirementSection(t *testing.T) {
	section := RequirementSection("3.8.constructor-category-order")
	if section != "3.8" {
		t.Fatalf("unexpected section %q", section)
	}
}
