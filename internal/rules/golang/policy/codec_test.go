package policy_test

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

	corepolicy "ciphera/tools/internal/policy"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
)

/* ------------------------------------------ Decoding ------------------------------------------ */

func TestDecodeConfigReadsGoPackConfig(t *testing.T) {
	t.Parallel()

	config, err := gopolicy.DecodeConfig(corepolicy.PackConfig{
		"local_import_prefixes": []any{"ciphera"},
		"parameters": map[string]any{
			"secret_names": []any{"token"},
		},
		"constructors": map[string]any{
			"parameter_order": []any{
				map[string]any{
					"name":               "repository",
					"type_name_suffixes": []any{"Repository"},
				},
			},
		},
		"domain_values": map[string]any{
			"required_constructors": map[string]any{
				"SessionKey": []any{"ParseSessionKey"},
			},
		},
		"architecture": map[string]any{
			"layers": []any{
				map[string]any{
					"name":          "core",
					"package_roots": []any{"internal/core"},
					"may_import":    []any{"core"},
				},
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
	got, err := gopolicy.DecodeConfig(gopolicy.EncodeConfig(expected))
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	if diff := cmp.Diff(expected, got, cmpopts.EquateEmpty()); diff != "" {
		t.Fatalf("config mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeConfigRejectsUnknownFields(t *testing.T) {
	t.Parallel()

	pack := gopolicy.EncodeConfig(baselineConfig())
	pack["surprise"] = true

	_, err := gopolicy.DecodeConfig(pack)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "packs.go.surprise") {
		t.Fatalf("unexpected error: %v", err)
	}
}

/* ----------------------------------------- Validation ----------------------------------------- */

func TestValidateConfigRejectsInvalidArchitecture(t *testing.T) {
	t.Parallel()

	config := baselineConfig()
	config.Architecture.Layers[0].AllowedLayers = []string{"missing"}

	err := gopolicy.ValidateConfig(config)
	if err == nil {
		t.Fatal("expected validation error")
	}

	if !strings.Contains(err.Error(), "unknown layer") {
		t.Fatalf("unexpected error: %v", err)
	}
}

/* ------------------------------------------ Fixtures ------------------------------------------ */

func baselineConfig() (config gopolicy.Config) {
	return gopolicy.Config{
		LocalImportPrefixes: []string{"ciphera"},
		Parameters: gopolicy.ParameterConfig{
			SecretNames: []string{"token"},
		},
		Constructors: gopolicy.ConstructorConfig{
			ParameterOrder: []gopolicy.ParameterGroup{
				{Name: "repository", TypeNameSuffixes: []string{"Repository"}},
			},
		},
		DomainValues: gopolicy.DomainValueConfig{
			RequiredConstructors: gopolicy.DomainValueConstructors{
				"SessionKey": []string{"ParseSessionKey"},
			},
		},
		Architecture: gopolicy.ArchitectureConfig{
			Layers: []gopolicy.ArchitectureLayer{
				{
					Name:          "core",
					PackageRoots:  []string{"internal/core"},
					AllowedLayers: []string{"core"},
				},
			},
		},
	}
}
