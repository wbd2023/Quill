package policy_test

import (
	"testing"

	"ciphera/tools/internal/policy"
)

func TestPathClassesLookupPatterns(t *testing.T) {
	classes := policy.PathClasses{
		"go_source": {"cmd/", "internal/"},
	}

	patterns := classes.LookupPatterns("go_source")
	requireEqual(t, []string{"cmd/", "internal/"}, patterns)

	patterns[0] = "mutated/"
	requireEqual(t, []string{"cmd/", "internal/"}, classes.LookupPatterns("go_source"))
}

func TestPathClassesLookupPatternsHandlesMissingClasses(t *testing.T) {
	cases := []struct {
		name    string
		classes policy.PathClasses
	}{
		{name: "nil classes", classes: nil},
		{
			name:    "unknown class",
			classes: policy.PathClasses{"markdown": {".md"}},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			requireEqual(t, []string(nil), test.classes.LookupPatterns("go_source"))
		})
	}
}
