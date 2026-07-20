package markers_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/markers"
)

const longLineRule = "allow-long-line"

func TestTextFormatsCanonicalMarker(t *testing.T) {
	marker := markers.Text(longLineRule)
	if marker != "style: allow-long-line" {
		t.Fatalf("unexpected marker %q", marker)
	}
}

func TestBecauseFormatsCanonicalMarkerWithReason(t *testing.T) {
	marker := markers.Because(longLineRule, "shell output")
	if marker != "style: allow-long-line because: shell output" {
		t.Fatalf("unexpected marker %q", marker)
	}
}

func TestBecauseIgnoresBlankReason(t *testing.T) {
	marker := markers.Because(longLineRule, "   ")
	if marker != "style: allow-long-line" {
		t.Fatalf("unexpected marker %q", marker)
	}
}
