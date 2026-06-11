package policy_test

import (
	"strings"
	"testing"

	"ciphera/tools/internal/checks/vocabulary/policy"
)

func TestValidateConfigRejectsIncompleteVocabularyPolicy(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		Go: policy.GoConfig{
			ForbiddenTypeSuffixes: []string{"Repository"},
		},
	}

	err := policy.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.vocabulary.go.preferred_type_suffix") {
		t.Fatalf("unexpected error: %v", err)
	}
}
