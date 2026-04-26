package styleguide

import "testing"

func TestParseStyleHeading(t *testing.T) {
	section, title, found := parseHeadingText("### 3.2 Context, resources, and concurrency")
	if !found {
		t.Fatal("expected heading to parse")
	}

	if section != "3.2" || title != "Context, resources, and concurrency" {
		t.Fatalf("unexpected heading parse: %q %q", section, title)
	}
}

func TestExtractStyleHeadingsParsesSectionHeadingsAtAnyMarkdownLevel(t *testing.T) {
	headings := extractHeadings(t, "## 3.2 Context, resources, and concurrency\n")
	if len(headings) != 1 {
		t.Fatalf("expected 1 heading, got %d", len(headings))
	}

	if headings[0].Section != "3.2" || headings[0].Title != "Context, resources, and concurrency" {
		t.Fatalf("unexpected heading parse: %+v", headings[0])
	}
}
