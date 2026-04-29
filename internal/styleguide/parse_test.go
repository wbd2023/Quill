package styleguide

import "testing"

/* ------------------------------------- Requirement Parsing ------------------------------------ */

func TestParseBuildsRequirementsFromMetadata(t *testing.T) {
	document := parseDocument(t, styleDocument(
		"### 3.2 Context, resources, and concurrency",
		"",
		"<!-- style: id=3.2.ctx-first -->",
		"* `ctx context.Context` MUST be the first parameter when present.",
	))

	requireDocument(t, document, Document{
		Headings: []Heading{
			{
				Section: "3.2",
				Title:   "Context, resources, and concurrency",
			},
		},
		Requirements: []Requirement{
			{
				ID:      "3.2.ctx-first",
				Section: "3.2",
				Text:    "ctx context.Context MUST be the first parameter when present.",
			},
		},
	})
}

func TestParseNormalisesRequirementText(t *testing.T) {
	requirements := parseRequirements(t, styleDocument(
		"### 1.1 Example",
		"",
		"<!-- style: id=1.1.example -->",
		"* **Important** requirements with `code`",
		"  MUST collapse formatting and wrapped lines.",
	))

	requireRequirements(t, requirements, []Requirement{
		{
			ID:      "1.1.example",
			Section: "1.1",
			Text:    "Important requirements with code MUST collapse formatting and wrapped lines.",
		},
	})
}

func TestParseIgnoresFencedCodeBlocks(t *testing.T) {
	contents := styleDocument(
		"### 1.2 Line length",
		"",
		"<!-- style: id=1.2.max-line-length -->",
		"* Maximum 100 characters per line.",
		"",
		"```go",
		"<!-- style: id=9.9.fake-requirement -->",
		"* This is example text, not a real requirement.",
		"```",
		"",
		"<!-- style: id=1.2.tabs-count-four -->",
		"* Tabs count as 4 columns for line-length checks.",
	)

	requireRequirements(t, parseRequirements(t, contents), []Requirement{
		{
			ID:      "1.2.max-line-length",
			Section: "1.2",
			Text:    "Maximum 100 characters per line.",
		},
		{
			ID:      "1.2.tabs-count-four",
			Section: "1.2",
			Text:    "Tabs count as 4 columns for line-length checks.",
		},
	})
}

/* ---------------------------------- Invalid Requirement Flow ---------------------------------- */

func TestParseRejectsInvalidRequirementFlow(t *testing.T) {
	cases := []struct {
		name     string
		contents string
		expected string
	}{
		{
			name: "requirement before heading",
			contents: styleDocument(
				"<!-- style: id=1.1.example -->",
				"* Requirements need a section.",
			),
			expected: "appears before any STYLE.md section heading",
		},
		{
			name: "requirement from another section",
			contents: styleDocument(
				"### 3.2 Context, resources, and concurrency",
				"",
				"<!-- style: id=3.3.ctx-first -->",
				"* This requirement belongs elsewhere.",
			),
			expected: "appears under section",
		},
		{
			name: "duplicate requirement",
			contents: styleDocument(
				"### 1.1 Example",
				"",
				"<!-- style: id=1.1.example -->",
				"* First copy.",
				"<!-- style: id=1.1.example -->",
				"* Second copy.",
			),
			expected: "duplicate STYLE.md requirement",
		},
		{
			name: "metadata before heading",
			contents: styleDocument(
				"<!-- style: id=1.1.example -->",
				"",
				"### 1.1 Example",
			),
			expected: "must be followed by a requirement list item",
		},
		{
			name: "metadata before blank requirement",
			contents: styleDocument(
				"### 1.1 Example",
				"",
				"<!-- style: id=1.1.example -->",
				"*",
			),
			expected: "must be followed by a requirement list item",
		},
		{
			name: "metadata before another metadata comment",
			contents: styleDocument(
				"### 1.1 Example",
				"",
				"<!-- style: id=1.1.example -->",
				"",
				"<!-- style: id=1.1.other -->",
			),
			expected: `metadata for "1.1.other" appears before metadata for "1.1.example" ` +
				`has a requirement list item`,
		},
		{
			name: "metadata before non-metadata html",
			contents: styleDocument(
				"### 1.1 Example",
				"",
				"<!-- style: id=1.1.example -->",
				"",
				"<!-- note: not a requirement -->",
			),
			expected: "must be followed by a requirement list item",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse([]byte(test.contents), testStyleGuideConfig())
			requireErrorContains(t, err, test.expected)
		})
	}
}

/* ----------------------------------------- Diagnostics ---------------------------------------- */

func TestParseReportsSourcePositions(t *testing.T) {
	contents := styleDocument(
		"### 1.1 Example",
		"",
		"<!-- style: id=1.1.example -->",
		"",
		"### 1.2 Next",
	)

	_, err := Parse([]byte(contents), testStyleGuideConfig())
	requireErrorContains(
		t,
		err,
		"STYLE.md:3:1: STYLE.md metadata for \"1.1.example\"",
	)
}
