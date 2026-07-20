package vocabularypolicy_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/wbd2023/Quill/internal/checks/vocabularypolicy"
	corepolicy "github.com/wbd2023/Quill/internal/policy"
)

func TestDecodeConfigReadsVocabularyPackConfig(t *testing.T) {
	t.Parallel()

	config, err := vocabularypolicy.DecodeConfig(corepolicy.PackConfig{
		"go": map[string]any{
			"type_suffixes": map[string]any{
				"Repository": []any{"Store"},
			},
			"identifier_suffixes": map[string]any{
				"Repository": []any{"Repo"},
				"Service":    []any{"Svc"},
			},
		},
		"bash": map[string]any{
			"variable_names": map[string]any{
				"named_constant": []any{"x"},
			},
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
			"type_suffixes": map[string]any{
				"Repository": []any{"Store"},
			},
			"surprise": true,
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
			TypeSuffixes: map[string][]string{
				"Repository": {"Store"},
			},
			IdentifierSuffixes: map[string][]string{
				"Repository": {"Repo"},
				"Service":    {"Svc"},
			},
		},
		Bash: vocabularypolicy.BashConfig{
			VariableNames: map[string][]string{
				"named_constant": {"x"},
			},
		},
	}
}
