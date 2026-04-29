package styleguide

import (
	"testing"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/text"
)

func TestExtractInlineTextCollapsesMarkdownText(t *testing.T) {
	source := []byte("**Important** text with `code`\nwrapped text.\n")
	tree := goldmark.DefaultParser().Parse(text.NewReader(source))

	got := extractInlineText(source, tree)
	expected := "Important text with code wrapped text."
	if got != expected {
		t.Fatalf("unexpected inline text\nexpected: %q\nactual:   %q", expected, got)
	}
}

func TestExtractInlineTextSeparatesMarkdownBlocks(t *testing.T) {
	source := []byte("* First paragraph.\n\n  Second paragraph.\n")
	tree := goldmark.DefaultParser().Parse(text.NewReader(source))

	got := extractInlineText(source, tree)
	expected := "First paragraph. Second paragraph."
	if got != expected {
		t.Fatalf("unexpected inline text\nexpected: %q\nactual:   %q", expected, got)
	}
}
