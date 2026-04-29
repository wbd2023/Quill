package styleguide

import "testing"

/* ------------------------------------------- Parsing ------------------------------------------ */

func TestParseBlockMetadataAcceptsSupportedForms(t *testing.T) {
	cases := []struct {
		name     string
		payload  string
		expected metadataFields
	}{
		{
			name:    "id only",
			payload: "id=1.1.example",
			expected: metadataFields{
				id: "1.1.example",
			},
		},
		{
			name: "review metadata",
			payload: "id=1.1.example\n" +
				"mode=review_only\n" +
				"reason=Review this manually with a wrapped\n" +
				"  explanation.",
			expected: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: "Review this manually with a wrapped explanation.",
			},
		},
		{
			name: "equals in wrapped reason",
			payload: "id=1.1.example\n" +
				"mode=review_only\n" +
				"reason=Review shell values with\n" +
				"  key=value examples.",
			expected: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: "Review shell values with key=value examples.",
			},
		},
		{
			name: "unindented wrapped reason",
			payload: "id=1.1.example\n" +
				"mode=review_only\n" +
				"reason=Review this manually with a wrapped\n" +
				"explanation.",
			expected: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: "Review this manually with a wrapped explanation.",
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fields, err := parseBlockMetadata(test.payload)
			if err != nil {
				t.Fatalf("parseBlockMetadata: %v", err)
			}
			requireMetadataFields(t, fields, test.expected)
		})
	}
}

func TestParseBlockMetadataRejectsMalformedInput(t *testing.T) {
	cases := []struct {
		name     string
		payload  string
		expected string
	}{
		{
			name:     "duplicate field",
			payload:  "id=1.1.example\nid=1.1.other",
			expected: `duplicate "id"`,
		},
		{
			name:     "unknown field",
			payload:  "id=1.1.example\nowner=style-team",
			expected: `unknown style metadata field "owner"`,
		},
		{
			name:     "unknown field with punctuation",
			payload:  "id=1.1.example\nfield-name=value",
			expected: `unknown style metadata field "field-name"`,
		},
		{
			name:     "empty id",
			payload:  "id=",
			expected: `style metadata field "id" must not be empty`,
		},
		{
			name: "empty mode",
			payload: "id=1.1.example\n" +
				"mode=\n" +
				"reason=Review this manually.",
			expected: `style metadata field "mode" must not be empty`,
		},
		{
			name: "continued id",
			payload: "id=1.1.example\n" +
				"  extra",
			expected: "malformed style metadata comment near",
		},
		{
			name: "continued mode",
			payload: "id=1.1.example\n" +
				"mode=review_only\n" +
				"  extra",
			expected: "malformed style metadata comment near",
		},
		{
			name: "empty reason",
			payload: "id=1.1.example\n" +
				"mode=review_only\n" +
				"reason=",
			expected: `style metadata field "reason" must not be empty`,
		},
		{
			name: "unknown field after reason",
			payload: "id=1.1.example\n" +
				"mode=review_only\n" +
				"reason=Review this manually.\n" +
				"owner=style-team",
			expected: `unknown style metadata field "owner"`,
		},
		{
			name:     "continuation before field",
			payload:  "continued text\nid=1.1.example",
			expected: "malformed style metadata comment near",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseBlockMetadata(test.payload)
			requireErrorContains(t, err, test.expected)
		})
	}
}
