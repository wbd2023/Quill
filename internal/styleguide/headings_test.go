package styleguide

import "testing"

func TestParseHeadingAcceptsNumberedSection(t *testing.T) {
	heading, found := parseHeading("3.2 Context, resources, and concurrency")
	if !found {
		t.Fatal("expected heading to parse")
	}

	requireHeading(t, heading, Heading{
		Section: "3.2",
		Title:   "Context, resources, and concurrency",
	})
}

func TestParseHeadingRejectsMalformedHeadings(t *testing.T) {
	cases := []struct {
		name string
		line string
	}{
		{name: "empty heading", line: ""},
		{name: "missing title", line: "3.2"},
		{name: "invalid section", line: "3.x Context"},
		{name: "not a numbered heading", line: "Context"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, found := parseHeading(test.line)
			if found {
				t.Fatal("expected malformed heading to be rejected")
			}
		})
	}
}

func TestParseReadsSectionHeadingsAtAnyMarkdownLevel(t *testing.T) {
	requireHeadings(
		t,
		parseHeadings(t, "## 3.2 Context, resources, and concurrency\n"),
		[]Heading{
			{
				Section: "3.2",
				Title:   "Context, resources, and concurrency",
			},
		},
	)
}
