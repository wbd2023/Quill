package policy_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"ciphera/tools/internal/checks/vocabulary/policy"
	corepolicy "ciphera/tools/internal/policy"
)

func TestDecodeConfigReadsVocabularyPackConfig(t *testing.T) {
	t.Parallel()

	config, err := policy.DecodeConfig(corepolicy.PackConfig{
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
	got, err := policy.DecodeConfig(policy.EncodeConfig(expected))
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeConfigRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	_, err := policy.DecodeConfig(corepolicy.PackConfig{
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

func baselineConfig() (config policy.Config) {
	return policy.Config{
		Go: policy.GoConfig{
			ForbiddenTypeSuffixes:       []string{"Repository"},
			PreferredTypeSuffix:         "Store",
			ForbiddenIdentifierSuffixes: []string{"Repository"},
			PreferredIdentifierSuffix:   "Store",
		},
		Bash: policy.BashConfig{
			ForbiddenVariableNames: []string{"x"},
			PreferredVariableName:  "named_constant",
		},
	}
}
