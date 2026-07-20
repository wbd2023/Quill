package policy_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/policy"
)

func TestPathRolesLookupPatterns(t *testing.T) {
	roles := policy.PathRoles{
		"go_source": {"cmd/", "internal/"},
	}

	patterns := roles.LookupPatterns("go_source")
	requireEqual(t, []string{"cmd/", "internal/"}, patterns)

	patterns[0] = "mutated/"
	requireEqual(t, []string{"cmd/", "internal/"}, roles.LookupPatterns("go_source"))
}

func TestPathRolesLookupPatternsHandlesMissingClasses(t *testing.T) {
	cases := []struct {
		name  string
		roles policy.PathRoles
	}{
		{name: "nil roles", roles: nil},
		{
			name:  "unknown role",
			roles: policy.PathRoles{"markdown": {".md"}},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			requireEqual(t, []string(nil), test.roles.LookupPatterns("go_source"))
		})
	}
}
