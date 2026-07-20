package projectpolicy_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/wbd2023/Quill/internal/checks/projectpolicy"
	corepolicy "github.com/wbd2023/Quill/internal/policy"
)

/* ------------------------------------------ Decoding ------------------------------------------ */

func TestDecodeConfigReadsProjectPackConfig(t *testing.T) {
	t.Parallel()

	config, err := projectpolicy.DecodeConfig(corepolicy.PackConfig{
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
	got, err := projectpolicy.DecodeConfig(projectpolicy.EncodeConfig(expected))
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	if diff := cmp.Diff(expected, got); diff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeConfigRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	pack := projectpolicy.EncodeConfig(baselineConfig())
	section := pack["commands"].(map[string]any)
	section["surprise"] = true

	_, err := projectpolicy.DecodeConfig(pack)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.project.commands.surprise") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDecodeConfigRejectsOldQualitySurfaceKey(t *testing.T) {
	t.Parallel()

	_, err := projectpolicy.DecodeConfig(corepolicy.PackConfig{
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

	err := projectpolicy.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "duplicate name") {
		t.Fatalf("unexpected error: %v", err)
	}
}

/* ------------------------------------------ Fixtures ------------------------------------------ */

func baselineConfig() (config projectpolicy.Config) {
	return projectpolicy.Config{
		Commands: projectpolicy.CommandsConfig{
			Runner: projectpolicy.CommandsRunnerMake,
			Make: projectpolicy.MakeConfig{
				Path: "mk/quality.mk",
				RequiredVariables: []projectpolicy.MakefileVariable{
					{Name: "LINT_ARGS", Value: "--strict"},
				},
				RequiredTargets: []projectpolicy.MakefileTarget{
					{Name: "lint", RecipeLine: "go test ./..."},
				},
			},
		},
	}
}
