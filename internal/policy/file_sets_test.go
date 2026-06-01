package policy_test

import (
	"testing"

	"ciphera/tools/internal/policy"
)

func TestFileSetsLookup(t *testing.T) {
	fileSets := policy.FileSets{
		{Name: "markdown", Include: policy.FileSetInclude{Extensions: []string{".md"}}},
	}

	fileSet, found := fileSets.Lookup("markdown")
	if !found {
		t.Fatalf("expected file set lookup to find markdown")
	}

	requireEqual(t, policy.FileSetConfig{
		Name: "markdown",
		Include: policy.FileSetInclude{
			Extensions: []string{".md"},
		},
	}, fileSet)

	_, found = fileSets.Lookup("missing")
	if found {
		t.Fatalf("expected missing file set lookup to fail")
	}
}
