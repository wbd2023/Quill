package styleguide

import (
	"testing"

	"ciphera/tools/internal/style"

	"github.com/google/go-cmp/cmp"
)

type metadataView struct {
	ID     string
	Review Review
}

/* ------------------------------------------ Comments ------------------------------------------ */

func TestParseMetadataCommentAcceptsSupportedForms(t *testing.T) {
	cases := []struct {
		name     string
		comment  string
		expected metadataFields
	}{
		{
			name: "inline metadata",
			comment: `<!-- style: id=1.1.example mode=review_only ` +
				`reason="Review this manually." -->`,
			expected: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: "Review this manually.",
			},
		},
		{
			name: "block metadata",
			comment: "<!-- style:\n" +
				"id=1.1.example\n" +
				"mode=review_only\n" +
				"reason=Review this manually with a wrapped\n" +
				"  explanation.\n" +
				"-->",
			expected: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: "Review this manually with a wrapped explanation.",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fields, found, err := parseMetadataComment(test.comment)
			if err != nil {
				t.Fatalf("parseMetadataComment: %v", err)
			}
			if !found {
				t.Fatal("expected metadata comment to parse")
			}

			requireMetadataFields(t, fields, test.expected)
		})
	}
}

/* ------------------------------------------ Payloads ------------------------------------------ */

func TestExtractMetadataPayloadAcceptsStyleComments(t *testing.T) {
	cases := []struct {
		name     string
		comment  string
		expected string
	}{
		{
			name:     "inline metadata",
			comment:  "<!-- style: id=1.1.example -->",
			expected: "id=1.1.example",
		},
		{
			name:     "block metadata",
			comment:  "<!-- style:\nid=1.1.example\n-->",
			expected: "id=1.1.example",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			payload, found, err := extractMetadataPayload(test.comment)
			if err != nil {
				t.Fatalf("extractMetadataPayload: %v", err)
			}
			if !found {
				t.Fatal("expected metadata payload to parse")
			}
			if payload != test.expected {
				t.Fatalf(
					"unexpected metadata payload\nexpected: %q\nactual:   %q",
					test.expected,
					payload,
				)
			}
		})
	}
}

func TestExtractMetadataPayloadIgnoresOtherComments(t *testing.T) {
	cases := []struct {
		name    string
		comment string
	}{
		{name: "other comment", comment: "<!-- note: id=1.1.example -->"},
		{name: "style substring", comment: "<!-- note: style: blah -->"},
		{name: "wrong metadata prefix", comment: "<!-- not-style: id=1.1.example -->"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, found, err := extractMetadataPayload(test.comment)
			if err != nil {
				t.Fatalf("extractMetadataPayload: %v", err)
			}
			if found {
				t.Fatal("expected non-style comment to be ignored")
			}
		})
	}
}

func TestExtractMetadataPayloadRejectsMalformedInput(t *testing.T) {
	cases := []struct {
		name    string
		comment string
	}{
		{name: "not html metadata", comment: "style: id=1.1.example"},
		{name: "unterminated html metadata", comment: "<!-- style: id=1.1.example"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, _, err := extractMetadataPayload(test.comment)
			requireErrorContains(t, err, "malformed style metadata comment")
		})
	}
}

/* ------------------------------------ Requirement Metadata ------------------------------------ */

func TestBuildRequirementMetadataAcceptsSupportedFields(t *testing.T) {
	cases := []struct {
		name     string
		fields   metadataFields
		expected metadataView
	}{
		{
			name: "id only",
			fields: metadataFields{
				id: "1.1.example",
			},
			expected: metadataView{
				ID: "1.1.example",
			},
		},
		{
			name: "review metadata",
			fields: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: "Review this manually.",
			},
			expected: metadataView{
				ID: "1.1.example",
				Review: Review{
					Only:   true,
					Reason: "Review this manually.",
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			metadata, err := buildRequirementMetadata(
				test.fields,
				style.SectionSlug,
			)
			if err != nil {
				t.Fatalf("buildRequirementMetadata: %v", err)
			}

			requireRequirementMetadata(t, metadata, test.expected)
		})
	}
}

func TestBuildRequirementMetadataRejectsInvalidFields(t *testing.T) {
	cases := []struct {
		name     string
		fields   metadataFields
		expected string
	}{
		{
			name: "missing reason",
			fields: metadataFields{
				id:   "1.1.example",
				mode: "review_only",
			},
			expected: "mode and reason must appear together",
		},
		{
			name: "missing mode",
			fields: metadataFields{
				id:     "1.1.example",
				reason: "Review this manually.",
			},
			expected: "mode and reason must appear together",
		},
		{
			name: "unsupported mode",
			fields: metadataFields{
				id:     "1.1.example",
				mode:   "manual",
				reason: "Review this manually.",
			},
			expected: "unsupported style metadata mode",
		},
		{
			name: "invalid id",
			fields: metadataFields{
				id:     "1.1.Bad",
				mode:   "review_only",
				reason: "Review this manually.",
			},
			expected: "invalid requirement id",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := buildRequirementMetadata(
				test.fields,
				style.SectionSlug,
			)
			requireErrorContains(t, err, test.expected)
		})
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func requireRequirementMetadata(
	t *testing.T,
	metadata requirementMetadata,
	expected metadataView,
) {
	t.Helper()

	actual := metadataView{
		ID:     metadata.id.String(),
		Review: metadata.review,
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Fatalf("unexpected requirement metadata (-expected +actual):\n%s", diff)
	}
}
