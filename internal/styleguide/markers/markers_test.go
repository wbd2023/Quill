package markers_test

import (
	"testing"

	"ciphera/tools/internal/styleguide/markers"
)

func TestParseOnlyMatchesCommentContexts(t *testing.T) {
	line := `const prefix = "style: "`

	marker := markers.Parse(line)
	if marker.Status != markers.Absent {
		t.Fatalf("string literal line should not be treated as an exception marker")
	}
}

func TestParseAcceptsSupportedCommentForms(t *testing.T) {
	testCases := []string{
		`echo "value" # style: allow-long-line because: shell output`,
		`// style: allow-non-ascii because: protocol sample`,
		`/* style: allow-long-line because: generated sample */`,
		` * style: allow-long-line because: generated sample */`,
	}

	for _, line := range testCases {
		marker := markers.Parse(line)
		if marker.Status != markers.Valid {
			t.Fatalf("expected valid marker for %q", line)
		}

		if marker.Rule == "" {
			t.Fatalf("expected parsed rule for %q", line)
		}
	}
}

func TestParseReportsInvalidMarkers(t *testing.T) {
	line := "// " + "style: allow-long-line because:"

	marker := markers.Parse(line)
	if marker.Status != markers.Invalid {
		t.Fatalf("expected invalid marker status, got %q", marker.Status)
	}
}

func TestTextWithReasonFormatsCanonicalMarker(t *testing.T) {
	marker := markers.TextWithReason(markers.LongLine, "shell output")
	if marker != "style: allow-long-line because: shell output" {
		t.Fatalf("unexpected marker %q", marker)
	}
}
