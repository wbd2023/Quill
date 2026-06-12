package vocabularypolicy_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"ciphera/tools/internal/checks/vocabularypolicy"
	corepolicy "ciphera/tools/internal/policy"
)

func TestDecodeConfigReadsVocabularyPackConfig(t *testing.T) {
	t.Parallel()

	config, err := vocabularypolicy.DecodeConfig(corepolicy.PackConfig{
		"go": map[string]any{
			"forbidden_type_suffixes":       []any{"Repository"},
			"preferred_type_suffix":         "Store",
			"forbidden_identifier_suffixes": []any{"Repository"},
			"preferred_identifier_suffix":   "Store",
		},
		"bash": map[string]any{
			"forbidden_variable_names": []any{"x"},
			"preferred_variable_name":  "named_constant",
		},
	})
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	if diff := cmp.Diff(baselineConfig(), config); diff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeConfigReadsEncodedConfig(t *testing.T) {
	t.Parallel()

	expected := baselineConfig()
	got, err := vocabularypolicy.DecodeConfig(vocabularypolicy.EncodeConfig(expected))
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeConfigRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	_, err := vocabularypolicy.DecodeConfig(corepolicy.PackConfig{
		"go": map[string]any{
			"preferred_type_suffix": "Store",
			"surprise":              true,
		},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.vocabulary.go.surprise") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func baselineConfig() (config vocabularypolicy.Config) {
	return vocabularypolicy.Config{
		Go: vocabularypolicy.GoConfig{
			ForbiddenTypeSuffixes:       []string{"Repository"},
			PreferredTypeSuffix:         "Store",
			ForbiddenIdentifierSuffixes: []string{"Repository"},
			PreferredIdentifierSuffix:   "Store",
		},
		Bash: vocabularypolicy.BashConfig{
			ForbiddenVariableNames: []string{"x"},
			PreferredVariableName:  "named_constant",
		},
	}
}
