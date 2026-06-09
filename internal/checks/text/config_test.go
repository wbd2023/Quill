package text_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"ciphera/tools/internal/checks/text"
	"ciphera/tools/internal/policy"
)

func TestDecodeConfigReadsTextPackConfig(t *testing.T) {
	t.Parallel()

	config, err := text.DecodeConfig(policy.PackConfig{
		"section_headers": map[string]any{
			"large_min_lines":  int64(100),
			"short_max_lines":  int64(79),
			"max_header_count": int64(6),
			"generic_names":    []any{"Check", "Checks"},
			"structural_names": []any{"Types", "Helpers"},
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
	got, err := text.DecodeConfig(text.EncodeConfig(expected))
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeConfigRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	pack := text.EncodeConfig(baselineConfig())
	section := pack["section_headers"].(map[string]any)
	section["required_min_lines"] = int64(100)

	_, err := text.DecodeConfig(pack)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.text.section_headers.required_min_lines") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateConfigRejectsInvalidSectionHeaderPolicy(t *testing.T) {
	t.Parallel()

	config := baselineConfig()
	config.SectionHeaders.ShortMaxLines = config.SectionHeaders.LargeMinLines

	err := text.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "short_max_lines") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func baselineConfig() (config text.Config) {
	return text.Config{
		SectionHeaders: text.SectionHeaderConfig{
			LargeMinLines:   100,
			ShortMaxLines:   79,
			MaxHeaderCount:  6,
			GenericNames:    []string{"Check", "Checks"},
			StructuralNames: []string{"Types", "Helpers"},
		},
	}
}
