package styleguide

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

/* ------------------------------------------ Scanning ------------------------------------------ */

func TestMarkdownScannerEmitsDocumentEvents(t *testing.T) {
	events := scanTestDocument(
		"### 1.1 Example",
		"",
		"<!-- style: id=1.1.example -->",
		"* Example MUST hold.",
		"",
		"Paragraph text.",
	)

	expected := []documentEvent{
		{
			kind:     eventHeading,
			location: position{line: 1, column: 1},
			heading:  Heading{Section: "1.1", Title: "Example"},
		},
		{
			kind:     eventHTMLBlock,
			location: position{line: 3, column: 1},
			text:     "<!-- style: id=1.1.example -->\n",
		},
		{
			kind:     eventListItem,
			location: position{line: 4, column: 1},
			text:     "Example MUST hold.",
		},
		{
			kind:     eventBoundary,
			location: position{line: 6, column: 1},
		},
	}
	requireDocumentEvents(t, events, expected)
}

func TestMarkdownScannerIgnoresBoundariesInsideListItems(t *testing.T) {
	events := scanTestDocument(
		"* Example MUST hold.",
		"",
		"  Additional explanation.",
		"",
		"Paragraph text.",
	)

	expected := []documentEvent{
		{
			kind:     eventListItem,
			location: position{line: 1, column: 1},
			text:     "Example MUST hold. Additional explanation.",
		},
		{
			kind:     eventBoundary,
			location: position{line: 5, column: 1},
		},
	}
	requireDocumentEvents(t, events, expected)
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func scanTestDocument(lines ...string) (events []documentEvent) {
	source := []byte(styleDocument(lines...))
	file := newSourceFile("STYLE.md", source)
	tree := goldmark.DefaultParser().Parse(text.NewReader(source))
	return scanMarkdown(tree, file)
}

func requireDocumentEvents(t *testing.T, events []documentEvent, expected []documentEvent) {
	t.Helper()

	if diff := cmp.Diff(
		expected,
		events,
		cmp.AllowUnexported(documentEvent{}, position{}),
	); diff != "" {
		t.Fatalf("unexpected document events (-expected +actual):\n%s", diff)
	}
}
