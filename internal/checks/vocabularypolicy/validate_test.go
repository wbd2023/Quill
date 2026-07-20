package vocabularypolicy_test

import (
	"strings"
	"testing"

	"github.com/wbd2023/Quill/internal/checks/vocabularypolicy"
)

func TestValidateConfigRejectsEmptyPreferredName(t *testing.T) {
	t.Parallel()

	config := vocabularypolicy.Config{
		Go: vocabularypolicy.GoConfig{
			TypeSuffixes: map[string][]string{
				"": {"Store"},
			},
		},
	}

	err := vocabularypolicy.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.vocabulary.go.type_suffixes") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateConfigRejectsShorthandMappedToTwoPreferred(t *testing.T) {
	t.Parallel()

	config := vocabularypolicy.Config{
		Go: vocabularypolicy.GoConfig{
			IdentifierSuffixes: map[string][]string{
				"Repository": {"Repo"},
				"Service":    {"Repo"},
			},
		},
	}

	err := vocabularypolicy.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), `shorthand "Repo"`) {
		t.Fatalf("unexpected error: %v", err)
	}
}
