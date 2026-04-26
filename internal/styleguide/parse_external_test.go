package styleguide_test

import (
	"testing"

	"ciphera/tools/internal/styleguide"
)

func TestParseExposesDocumentModelThroughPublicAPI(t *testing.T) {
	document, err := styleguide.Parse(
		[]byte(
			"### 1.1 Example\n\n"+
				"* [1.1.example] Public parsing should expose requirements.\n",
		),
		styleguide.Config{RequirementIDFormat: styleguide.RequirementIDFormatSectionSlug},
	)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	if len(document.Headings) != 1 || len(document.Requirements) != 1 {
		t.Fatalf("unexpected document shape: %+v", document)
	}

	if document.Requirements[0].ID != "1.1.example" {
		t.Fatalf("unexpected requirement ID %q", document.Requirements[0].ID)
	}
}
