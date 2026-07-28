package engine

import (
	"context"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/policy"
	"github.com/wbd2023/quill/internal/styleguide"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

// undocumentedRequirementID is syntactically valid (parses as section "9.9", slug
// "never-documented") but is absent from this repository's STYLE.md, so profile loading accepts it
// while requirement-binding validation must reject it.
const undocumentedRequirementID = "9.9.never-documented"

/* ------------------------------------ Requirement Bindings ------------------------------------ */

func TestValidateRequirementBindingsAcceptsDocumentedIDs(t *testing.T) {
	t.Parallel()

	config := policy.Profile{
		StyleGuide: policy.StyleGuideConfig{Path: "STYLE.md"},
		Rules: []policy.RuleBinding{
			{RuleID: "go/lint", RequirementIDs: []string{"3.1.foo", "3.2.bar"}},
		},
	}
	document := styleguide.Document{
		Requirements: []styleguide.Requirement{
			{ID: "3.1.foo"},
			{ID: "3.2.bar"},
			{ID: "4.5.unused"},
		},
	}

	if err := validateRequirementBindings(config, document); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRequirementBindingsRejectsUndocumentedID(t *testing.T) {
	t.Parallel()

	config := policy.Profile{
		StyleGuide: policy.StyleGuideConfig{Path: "STYLE.md"},
		Rules: []policy.RuleBinding{
			{RuleID: "go/lint", RequirementIDs: []string{"3.1.foo", "9.9.missing"}},
		},
	}
	document := styleguide.Document{
		Requirements: []styleguide.Requirement{{ID: "3.1.foo"}},
	}

	err := validateRequirementBindings(config, document)
	if err == nil {
		t.Fatalf("expected an error for an undocumented requirement id")
	}
	if !strings.Contains(err.Error(), `"9.9.missing"`) {
		t.Fatalf("error = %q, want it to name the undocumented requirement id", err.Error())
	}
	if !strings.Contains(err.Error(), "STYLE.md") {
		t.Fatalf("error = %q, want it to name the style guide path", err.Error())
	}
}

func TestValidateRequirementBindingsReportsAllUndocumentedIDsSorted(t *testing.T) {
	t.Parallel()

	config := policy.Profile{
		StyleGuide: policy.StyleGuideConfig{Path: "STYLE.md"},
		Rules: []policy.RuleBinding{
			{RuleID: "a", RequirementIDs: []string{"9.9.zeta"}},
			{RuleID: "b", RequirementIDs: []string{"9.9.alpha", "9.9.zeta"}},
		},
	}
	document := styleguide.Document{}

	err := validateRequirementBindings(config, document)
	if err == nil {
		t.Fatalf("expected an error for undocumented requirement ids")
	}

	message := err.Error()
	if !strings.Contains(message, `"9.9.alpha"`) || !strings.Contains(message, `"9.9.zeta"`) {
		t.Fatalf("error = %q, want both undocumented ids named", message)
	}
	if strings.Index(message, `"9.9.alpha"`) > strings.Index(message, `"9.9.zeta"`) {
		t.Fatalf("error = %q, want undocumented ids reported in sorted order", message)
	}
}

func TestCheckFailsOnUndocumentedRequirementBeforeToolInspection(t *testing.T) {
	t.Parallel()

	runner := &recordingRunner{}
	engine, err := New(undocumentedRequirementRepository(t), WithCommandRunner(runner))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = engine.Check(context.Background(), CheckOptions{})
	if err == nil {
		t.Fatalf("Check succeeded for an undocumented requirement id; want an error")
	}
	if !strings.Contains(err.Error(), undocumentedRequirementID) {
		t.Fatalf("Check error = %q, want it to name %q", err.Error(), undocumentedRequirementID)
	}
	if runner.calls != 0 {
		t.Fatalf("Check inspected tools %d time(s) before failing; want 0", runner.calls)
	}
}

func TestCoverageUsesSharedPreparationPipeline(t *testing.T) {
	t.Parallel()

	engine, err := New(undocumentedRequirementRepository(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = engine.Coverage(context.Background())
	if err == nil {
		t.Fatalf("Coverage succeeded for an undocumented requirement id; want an error")
	}
	if !strings.Contains(err.Error(), undocumentedRequirementID) {
		t.Fatalf("Coverage error = %q, want it to name %q", err.Error(), undocumentedRequirementID)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// undocumentedRequirementRepository writes a temporary repository whose first profile rule
// binds a syntactically valid requirement id absent from the copied STYLE.md. Every operation
// must reject it through the shared preparation pipeline before reaching tools or rules.
func undocumentedRequirementRepository(t *testing.T) (root string) {
	t.Helper()

	root = t.TempDir()
	config := profiles.Self(t)
	config.Rules[0].RequirementIDs = []string{undocumentedRequirementID}
	profiles.Write(t, root, config)
	return root
}
