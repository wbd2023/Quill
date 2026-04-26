package styleguide

import "testing"

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
