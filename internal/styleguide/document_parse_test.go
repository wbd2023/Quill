package styleguide

import "testing"

/* -------------------------------------- Document Parsing -------------------------------------- */

func TestParseStyleHeading(t *testing.T) {
	section, title, found := parseHeading("### 3.2 Context, resources, and concurrency")
	if !found {
		t.Fatal("expected heading to parse")
	}

	if section != "3.2" || title != "Context, resources, and concurrency" {
		t.Fatalf("unexpected heading parse: %q %q", section, title)
	}
}

func TestExtractStyleHeadingsParsesSectionHeadingsAtAnyMarkdownLevel(t *testing.T) {
	headings := extractHeadings("## 3.2 Context, resources, and concurrency\n")
	if len(headings) != 1 {
		t.Fatalf("expected 1 heading, got %d", len(headings))
	}

	if headings[0].Section != "3.2" || headings[0].Title != "Context, resources, and concurrency" {
		t.Fatalf("unexpected heading parse: %+v", headings[0])
	}
}

func TestParseStyleRequirementID(t *testing.T) {
	requirementID, found := parseRequirementID(
		"* `[3.2.ctx-first]` `ctx context.Context` MUST be the first parameter when present.",
	)
	if !found {
		t.Fatal("expected requirement ID to parse")
	}

	if requirementID != "3.2.ctx-first" {
		t.Fatalf("unexpected requirement ID %q", requirementID)
	}
}

func TestParseStyleRequirement(t *testing.T) {
	requirementID, requirementText, found := parseRequirement(
		"* `[3.2.ctx-first]` `ctx context.Context` MUST be the first parameter when present.",
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

func TestExtractStyleRequirementsIgnoresFencedCodeBlocks(t *testing.T) {
	contents := "### 1.2 Line length\n\n" +
		"* [1.2.max-line-length] Maximum 100 characters per line.\n\n" +
		"```go\n" +
		"* [9.9.fake-requirement] This is example text, not a real requirement.\n" +
		"```\n\n" +
		"* [1.2.tabs-count-four] Tabs count as 4 columns for line-length checks.\n"

	requirements := extractRequirements(contents)
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

func TestParseStyleMetadataComment(t *testing.T) {
	metadata, found, err := parseMetadataComment(
		`<!-- style: id=1.1.example mode=review_only reason="Review this manually." -->`,
		RequirementIDFormatSectionSlug,
	)
	if err != nil {
		t.Fatalf("parseMetadataComment: %v", err)
	}
	if !found {
		t.Fatal("expected metadata comment to parse")
	}
	if metadata.Mode != VerificationReviewOnly || metadata.Reason != "Review this manually." {
		t.Fatalf("unexpected metadata parse: %+v", metadata)
	}
}

func TestParseStyleMetadataCommentBlock(t *testing.T) {
	metadata, found, err := parseMetadataComment(
		"<!-- style:\n"+
			"id=1.1.example\n"+
			"mode=review_only\n"+
			"reason=Review this manually with a wrapped\n"+
			"  explanation.\n"+
			"-->",
		RequirementIDFormatSectionSlug,
	)
	if err != nil {
		t.Fatalf("parseMetadataComment block: %v", err)
	}
	if !found {
		t.Fatal("expected metadata block to parse")
	}
	if metadata.Mode != VerificationReviewOnly {
		t.Fatalf("unexpected metadata mode: %+v", metadata)
	}
	if metadata.Reason != "Review this manually with a wrapped explanation." {
		t.Fatalf("unexpected metadata reason %q", metadata.Reason)
	}
}

func TestCompileStyleDocumentRejectsMalformedMetadata(t *testing.T) {
	_, err := compileDocument([]byte(
		"### 1.1 Example\n\n<!-- style: id=1.1.example mode=review_only -->\n* Example.\n",
	), testStyleGuideConfig())
	if err == nil {
		t.Fatal("expected malformed metadata error")
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

func TestRequirementSection(t *testing.T) {
	section := RequirementSection("3.8.constructor-category-order")
	if section != "3.8" {
		t.Fatalf("unexpected section %q", section)
	}
}
