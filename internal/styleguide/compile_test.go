package styleguide

import (
	"testing"

	"ciphera/tools/internal/requirementid"
)

func TestDocumentCompilerBuildsDocumentFromEvents(t *testing.T) {
	compiler := newTestDocumentCompiler()

	document, err := compiler.compile([]documentEvent{
		newHeadingEvent(
			position{line: 1, column: 1},
			Heading{
				Section: "1.1",
				Title:   "Example",
			},
		),
		newHTMLBlockEvent(position{line: 3, column: 1}, "<!-- style: id=1.1.example -->"),
		newListItemEvent(position{line: 4, column: 1}, "Example MUST hold."),
	})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}

	requireDocument(t, document, Document{
		Headings: []Heading{
			{
				Section: "1.1",
				Title:   "Example",
			},
		},
		Requirements: []Requirement{
			{
				ID:      "1.1.example",
				Section: "1.1",
				Text:    "Example MUST hold.",
			},
		},
	})
}

func TestDocumentCompilerRejectsInterruptedMetadata(t *testing.T) {
	compiler := newTestDocumentCompiler()

	_, err := compiler.compile([]documentEvent{
		newHeadingEvent(
			position{line: 1, column: 1},
			Heading{
				Section: "1.1",
				Title:   "Example",
			},
		),
		newHTMLBlockEvent(position{line: 3, column: 1}, "<!-- style: id=1.1.example -->"),
		newBoundaryEvent(position{line: 5, column: 1}),
	})
	requireErrorContains(
		t,
		err,
		"STYLE.md:3:1: STYLE.md metadata for \"1.1.example\"",
	)
}

func TestDocumentCompilerRejectsUnknownEvents(t *testing.T) {
	compiler := newTestDocumentCompiler()

	_, err := compiler.compile([]documentEvent{{kind: eventKind(99)}})
	requireErrorContains(t, err, "unknown styleguide document event 99")
}

func newTestDocumentCompiler() (compiler documentCompiler) {
	return newDocumentCompiler(newSourceFile("STYLE.md", nil), requirementid.SectionSlug)
}
