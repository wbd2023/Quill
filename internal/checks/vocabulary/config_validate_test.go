package vocabulary_test

import (
	"strings"
	"testing"

	"ciphera/tools/internal/checks/vocabulary"
)

func TestValidateConfigRejectsIncompleteVocabularyPolicy(t *testing.T) {
	t.Parallel()

	config := vocabulary.Config{
		Go: vocabulary.GoConfig{
			ForbiddenTypeSuffixes: []string{"Repository"},
		},
	}

	err := vocabulary.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.vocabulary.go.preferred_type_suffix") {
		t.Fatalf("unexpected error: %v", err)
	}
}
