package styleguide

import "testing"

/* ------------------------------------------- Parsing ------------------------------------------ */

func TestParseInlineMetadataAcceptsSupportedForms(t *testing.T) {
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
			payload: `id=1.1.example mode=review_only ` +
				`reason="Review this manually."`,
			expected: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: "Review this manually.",
			},
		},
		{
			name: "escaped quoted reason",
			payload: `id=1.1.example   mode=review_only ` +
				`reason="Review \"manual\" checks."`,
			expected: metadataFields{
				id:     "1.1.example",
				mode:   "review_only",
				reason: `Review "manual" checks.`,
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			fields, err := parseInlineMetadata(test.payload)
			if err != nil {
				t.Fatalf("parseInlineMetadata: %v", err)
			}
			requireMetadataFields(t, fields, test.expected)
		})
	}
}

func TestParseInlineMetadataRejectsMalformedInput(t *testing.T) {
	cases := []struct {
		name     string
		payload  string
		expected string
	}{
		{
			name:     "missing id",
			payload:  `mode=review_only reason="Review this manually."`,
			expected: "malformed style metadata comment",
		},
		{
			name:     "missing mode",
			payload:  `id=1.1.example reason="Review this manually."`,
			expected: "malformed style metadata comment",
		},
		{
			name:     "reason before mode",
			payload:  `id=1.1.example reason="Review this manually." mode=review_only`,
			expected: "malformed style metadata comment",
		},
		{
			name:     "unknown field",
			payload:  "id=1.1.example owner=style-team",
			expected: `unknown style metadata field "owner"`,
		},
		{
			name:     "empty id",
			payload:  "id=",
			expected: `style metadata field "id" must not be empty`,
		},
		{
			name:     "empty mode",
			payload:  `id=1.1.example mode= reason="Review this manually."`,
			expected: `style metadata field "mode" must not be empty`,
		},
		{
			name:     "empty reason",
			payload:  `id=1.1.example mode=review_only reason=""`,
			expected: `style metadata field "reason" must not be empty`,
		},
		{
			name:     "unterminated reason",
			payload:  `id=1.1.example mode=review_only reason="Review this manually.`,
			expected: "malformed style metadata comment",
		},
		{
			name: "trailing text after reason",
			payload: `id=1.1.example mode=review_only ` +
				`reason="Review this manually." trailing`,
			expected: "malformed style metadata comment",
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := parseInlineMetadata(test.payload)
			requireErrorContains(t, err, test.expected)
		})
	}
}
