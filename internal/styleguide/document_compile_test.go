package styleguide

import "testing"

func TestExtractStyleRequirementsIgnoresFencedCodeBlocks(t *testing.T) {
	contents := "### 1.2 Line length\n\n" +
		"* [1.2.max-line-length] Maximum 100 characters per line.\n\n" +
		"```go\n" +
		"* [9.9.fake-requirement] This is example text, not a real requirement.\n" +
		"```\n\n" +
		"* [1.2.tabs-count-four] Tabs count as 4 columns for line-length checks.\n"

	requirements := extractRequirements(t, contents)
	if len(requirements) != 2 {
		t.Fatalf("expected 2 requirements, got %d", len(requirements))
	}

	if requirements[0].ID != "1.2.max-line-length" {
		t.Fatalf("unexpected first requirement %q", requirements[0].ID)
	}

	if requirements[1].ID != "1.2.tabs-count-four" {
		t.Fatalf("unexpected second requirement %q", requirements[1].ID)
	}
}

func TestExtractStyleRequirementsIgnoresRequirementIDsFromOtherSections(t *testing.T) {
	_, err := compileDocument([]byte(
		"### 3.2 Context, resources, and concurrency\n\n"+
			"* [3.3.ctx-first] This requirement ID belongs to another section.\n",
	), testStyleGuideConfig())
	if err == nil {
		t.Fatal("expected mismatched-section requirement to be rejected")
	}
}

func TestCompileStyleDocumentReadsHiddenRequirementIDs(t *testing.T) {
	document, err := compileDocument([]byte(
		"### 3.2 Context, resources, and concurrency\n\n"+
			"<!-- style: id=3.2.ctx-first -->\n"+
			"* `ctx context.Context` MUST be the first parameter when present.\n",
	), testStyleGuideConfig())
	if err != nil {
		t.Fatalf("compileDocument: %v", err)
	}

	if len(document.Requirements) != 1 {
		t.Fatalf("expected 1 requirement, got %d", len(document.Requirements))
	}

	requirement := document.Requirements[0]
	if requirement.ID != "3.2.ctx-first" {
		t.Fatalf("unexpected requirement id %q", requirement.ID)
	}

	if requirement.Text != "ctx context.Context MUST be the first parameter when present." {
		t.Fatalf("unexpected requirement text %q", requirement.Text)
	}
}
