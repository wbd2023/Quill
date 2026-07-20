package markers_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/markers"
)

func TestParseOnlyMatchesCommentContexts(t *testing.T) {
	line := `const prefix = "style: "`

	marker := markers.Parse(line)
	if marker.Status != markers.StatusAbsent {
		t.Fatalf("string literal line should not be treated as an exception marker")
	}
}

func TestParseAcceptsSupportedCommentForms(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{name: "hash comment", line: `echo "value" # style: allow-long-line because: shell output`},
		{name: "slash comment", line: `// style: allow-non-ascii because: protocol sample`},
		{name: "block comment", line: `/* style: allow-long-line because: generated sample */`},
		{
			name: "block continuation",
			line: ` * style: allow-long-line because: generated sample */`,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			marker := markers.Parse(test.line)
			if marker.Status != markers.StatusValid {
				t.Fatalf("expected valid marker for %q", test.line)
			}

			if marker.Rule == "" {
				t.Fatalf("expected parsed rule for %q", test.line)
			}
		})
	}
}

func TestParseReportsInvalidMarkers(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{name: "missing reason", line: "// style: allow-long-line because:"},
		{name: "empty rule segment", line: "// style: allow--line"},
		{name: "trailing rule hyphen", line: "// style: allow-long-line-"},
		{name: "unsupported rule character", line: "// style: allow-long_line"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			marker := markers.Parse(test.line)
			if marker.Status != markers.StatusInvalid {
				t.Fatalf("expected invalid marker status, got %q", marker.Status)
			}
		})
	}
}
