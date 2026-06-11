package policy_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"ciphera/tools/internal/checks/project/policy"
	corepolicy "ciphera/tools/internal/policy"
)

/* ------------------------------------------ Decoding ------------------------------------------ */

func TestDecodeConfigReadsProjectPackConfig(t *testing.T) {
	t.Parallel()

	config, err := policy.DecodeConfig(corepolicy.PackConfig{
		"commands": map[string]any{
			"runner": "make",
			"path":   "mk/quality.mk",
			"required_variables": []any{
				map[string]any{"name": "LINT_ARGS", "value": "--strict"},
			},
			"required_targets": []any{
				map[string]any{"name": "lint", "recipe_line": "go test ./..."},
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

	pack := policy.EncodeConfig(baselineConfig())
	section := pack["commands"].(map[string]any)
	section["surprise"] = true

	_, err := policy.DecodeConfig(pack)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.project.commands.surprise") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDecodeConfigRejectsOldQualitySurfaceKey(t *testing.T) {
	t.Parallel()

	_, err := policy.DecodeConfig(corepolicy.PackConfig{
		"quality_surface": map[string]any{},
	})
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.project.quality_surface") {
		t.Fatalf("unexpected error: %v", err)
	}
}

/* ----------------------------------------- Validation ----------------------------------------- */

func TestValidateConfigRejectsDuplicateTargets(t *testing.T) {
	t.Parallel()

	config := baselineConfig()
	target := config.Commands.Make.RequiredTargets[0]
	config.Commands.Make.RequiredTargets = append(
		config.Commands.Make.RequiredTargets,
		target,
	)

	err := policy.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "duplicate name") {
		t.Fatalf("unexpected error: %v", err)
	}
}

/* ------------------------------------------ Fixtures ------------------------------------------ */

func baselineConfig() (config policy.Config) {
	return policy.Config{
		Commands: policy.CommandsConfig{
			Runner: policy.CommandsRunnerMake,
			Make: policy.MakeConfig{
				Path: "mk/quality.mk",
				RequiredVariables: []policy.MakefileVariable{
					{Name: "LINT_ARGS", Value: "--strict"},
				},
				RequiredTargets: []policy.MakefileTarget{
					{Name: "lint", RecipeLine: "go test ./..."},
				},
			},
		},
	}
}
