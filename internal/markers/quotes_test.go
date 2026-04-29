package markers_test

import (
	"testing"

	"ciphera/tools/internal/markers"
)

func TestParseIgnoresMarkersInsideQuotedText(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{name: "double quoted marker", line: `const text = "// style: allow-long-line"`},
		{name: "single quoted marker", line: `echo '# style: allow-long-line'`},
		{name: "backtick quoted marker", line: "const text = `/* style: allow-long-line */`"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			marker := markers.Parse(test.line)
			if marker.Status != markers.StatusAbsent {
				t.Fatalf("expected marker inside quoted text to be absent for %q", test.line)
			}
		})
	}
}

func TestParseAcceptsMarkersAfterQuotedText(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{
			name: "after double quoted text",
			line: `fmt.Println("// not a marker") // style: allow-long-line`,
		},
		{name: "after single quoted text", line: `echo '# not a marker' # style: allow-long-line`},
		{
			name: "after backtick quoted text",
			line: "const text = `/* not a marker */` // style: allow-long-line",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			marker := markers.Parse(test.line)
			if marker.Status != markers.StatusValid {
				t.Fatalf("expected marker after quoted text for %q", test.line)
			}
		})
	}
}

func TestParseHonoursEscapedQuotes(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{name: "escaped double quote", line: `const text = "\" // style: allow-long-line"`},
		{name: "escaped single quote", line: `echo 'value \' # style: allow-long-line'`},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			marker := markers.Parse(test.line)
			if marker.Status != markers.StatusAbsent {
				t.Fatalf("expected escaped quote marker to remain quoted for %q", test.line)
			}
		})
	}
}
