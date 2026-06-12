package vocabularypolicy_test

import (
	"strings"
	"testing"

	"ciphera/tools/internal/checks/vocabularypolicy"
)

func TestValidateConfigRejectsIncompleteVocabularyPolicy(t *testing.T) {
	t.Parallel()

	config := vocabularypolicy.Config{
		Go: vocabularypolicy.GoConfig{
			ForbiddenTypeSuffixes: []string{"Repository"},
		},
	}

	err := vocabularypolicy.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.vocabulary.go.preferred_type_suffix") {
		t.Fatalf("unexpected error: %v", err)
	}
}
