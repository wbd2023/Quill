package styleguide

import "testing"

func TestParseExceptionMarkerOnlyMatchesCommentContexts(t *testing.T) {
	line := `const exceptionPrefix = "style: "`

	_, _, found, valid := ParseExceptionMarker(line)
	if found || valid {
		t.Fatalf("string literal line should not be treated as an exception marker")
	}
}

func TestParseExceptionMarkerAcceptsShellAndGoComments(t *testing.T) {
	testCases := []string{
		`echo "value" # style: allow-long-line because: shell output`,
		`// style: allow-non-ascii because: protocol sample`,
	}

	for _, line := range testCases {
		rule, _, found, valid := ParseExceptionMarker(line)
		if !found || !valid {
			t.Fatalf("expected valid marker for %q", line)
		}

		if rule == "" {
			t.Fatalf("expected parsed rule for %q", line)
		}
	}
}
