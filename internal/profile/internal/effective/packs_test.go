package effective_test

import (
	"testing"

	textrules "ciphera/tools/internal/checks/text"
	vocabularyrules "ciphera/tools/internal/checks/vocabulary"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/pack/shipped/bash"
	"ciphera/tools/internal/pack/shipped/markdown"
	"ciphera/tools/internal/pack/shipped/text"
	"ciphera/tools/internal/pack/shipped/vocabulary"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/effective"
)

/* ----------------------------------------- Validation ----------------------------------------- */

func TestResolvePacksRejectsMissingRequiredConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{vocabulary.PackID},
	}
	registry := registryFor(t, config)

	_, err := effective.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.vocabulary")
}

func TestResolvePacksRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{vocabulary.PackID},
		PackConfigs: policy.PackConfigs{
			vocabulary.PackID: vocabularyrules.EncodeConfig(vocabularyrules.Config{
				Go: vocabularyrules.GoConfig{ForbiddenTypeSuffixes: []string{"Repository"}},
			}),
		},
	}
	registry := registryFor(t, config)

	_, err := effective.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.vocabulary.go.preferred_type_suffix")
}

func TestResolvePacksRejectsInvalidTextConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{text.PackID},
		PackConfigs: policy.PackConfigs{
			text.PackID: textrules.EncodeConfig(textrules.Config{}),
		},
	}
	registry := registryFor(t, config)

	_, err := effective.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.text.section_headers")
}

func TestResolvePacksRejectsUnsupportedConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{markdown.PackID},
		PackConfigs: policy.PackConfigs{
			markdown.PackID: {"unknown": true},
		},
	}
	registry := registryFor(t, config)

	_, err := effective.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.markdown config is not supported")
}

func TestResolvePacksAcceptsPacksWithoutConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{markdown.PackID},
	}
	registry := registryFor(t, config)

	if _, err := effective.ResolvePacks(config, registry.Packs()); err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}
}

/* -------------------------------------- Default File Sets ------------------------------------- */

func TestResolvePacksAppliesPackDefaultFileSets(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{bash.PackID},
	}
	registry := registryFor(t, config)

	resolved, err := effective.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}

	fileSet, found := resolved.FileSets.Lookup("bash")
	if !found {
		t.Fatal("expected bash file set")
	}

	if len(fileSet.Include.Extensions) != 1 || fileSet.Include.Extensions[0] != ".sh" {
		t.Fatalf("bash extensions = %#v", fileSet.Include.Extensions)
	}
}

func TestResolvePacksLetsProfileFileSetsOverridePackDefaults(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{bash.PackID},
		FileSets: policy.FileSets{
			{Name: "bash", Include: policy.FileSetInclude{Extensions: []string{".bash"}}},
		},
	}
	registry := registryFor(t, config)

	resolved, err := effective.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}

	fileSet, found := resolved.FileSets.Lookup("bash")
	if !found {
		t.Fatal("expected bash file set")
	}

	if len(fileSet.Include.Extensions) != 1 || fileSet.Include.Extensions[0] != ".bash" {
		t.Fatalf("bash extensions = %#v", fileSet.Include.Extensions)
	}
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func registryFor(t *testing.T, config policy.Config) (registry pack.Registry) {
	t.Helper()

	registry, err := shipped.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	return registry
}
