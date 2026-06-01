package styleguide_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ciphera/tools/internal/requirementid"
	"ciphera/tools/internal/styleguide"

	"github.com/google/go-cmp/cmp"
)

/* -------------------------------------------- Parse ------------------------------------------- */

func TestParseExposesDocumentModelThroughPublicAPI(t *testing.T) {
	document, err := styleguide.Parse(
		[]byte(styleDocument(
			"### 1.1 Example",
			"",
			"<!-- style: id=1.1.example -->",
			"* Public parsing should expose requirements.",
		)),
		styleguide.Config{IDScheme: requirementid.SectionSlug},
	)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	requireDocument(t, document, styleguide.Document{
		Headings: []styleguide.Heading{
			{
				Section: "1.1",
				Title:   "Example",
			},
		},
		Requirements: []styleguide.Requirement{
			{
				ID:      "1.1.example",
				Section: "1.1",
				Text:    "Public parsing should expose requirements.",
			},
		},
	})
}

func TestParseExposesReviewMetadataThroughPublicAPI(t *testing.T) {
	document, err := styleguide.Parse(
		[]byte(styleDocument(
			"### 1.1 Example",
			"",
			`<!-- style: id=1.1.example mode=review_only reason="Review this manually." -->`,
			"* Public parsing should expose review metadata.",
		)),
		styleguide.Config{IDScheme: requirementid.SectionSlug},
	)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}

	requireDocument(t, document, styleguide.Document{
		Headings: []styleguide.Heading{
			{
				Section: "1.1",
				Title:   "Example",
			},
		},
		Requirements: []styleguide.Requirement{
			{
				ID:      "1.1.example",
				Section: "1.1",
				Text:    "Public parsing should expose review metadata.",
				Review: styleguide.Review{
					Only:   true,
					Reason: "Review this manually.",
				},
			},
		},
	})
}

func TestParseRejectsInvalidConfig(t *testing.T) {
	cases := []struct {
		name     string
		config   styleguide.Config
		expected string
	}{
		{
			name:     "missing requirement id scheme",
			config:   styleguide.Config{},
			expected: "requirement id scheme must not be empty",
		},
		{
			name: "unsupported requirement id scheme",
			config: styleguide.Config{
				IDScheme: "section",
			},
			expected: "unsupported styleguide requirement id scheme",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := styleguide.Parse(nil, test.config)
			requireErrorContains(t, err, test.expected)
		})
	}
}

/* -------------------------------------------- Load -------------------------------------------- */

func TestLoadReadsConfiguredStyleGuide(t *testing.T) {
	root := t.TempDir()
	filename := filepath.Join("docs", "STYLE.md")
	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}

	contents := styleDocument(
		"### 1.1 Example",
		"",
		"<!-- style: id=1.1.example -->",
		"* Loading from disk should expose requirements.",
	)
	if err := os.WriteFile(filepath.Join(root, filename), []byte(contents), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	document, err := styleguide.Load(root, styleguide.Config{
		Filename: filename,
		IDScheme: requirementid.SectionSlug,
	})
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	requireDocument(t, document, styleguide.Document{
		Headings: []styleguide.Heading{
			{
				Section: "1.1",
				Title:   "Example",
			},
		},
		Requirements: []styleguide.Requirement{
			{
				ID:      "1.1.example",
				Section: "1.1",
				Text:    "Loading from disk should expose requirements.",
			},
		},
	})
}

func TestLoadRejectsMissingFilename(t *testing.T) {
	_, err := styleguide.Load(t.TempDir(), styleguide.Config{
		IDScheme: requirementid.SectionSlug,
	})
	requireErrorContains(t, err, "styleguide filename must not be empty")
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func styleDocument(lines ...string) (document string) {
	return strings.Join(lines, "\n") + "\n"
}

func requireDocument(t *testing.T, document styleguide.Document, expected styleguide.Document) {
	t.Helper()

	if diff := cmp.Diff(expected, document); diff != "" {
		t.Fatalf("unexpected document (-expected +actual):\n%s", diff)
	}
}

func requireErrorContains(t *testing.T, err error, expected string) {
	t.Helper()

	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected error containing %q, got %v", expected, err)
	}
}
