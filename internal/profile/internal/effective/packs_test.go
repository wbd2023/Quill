package effective_test

import (
	"testing"

	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/rules/vocabulary"
)

/* ----------------------------------------- Validation ----------------------------------------- */

func TestResolvePacksRejectsMissingRequiredConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{builtin.PackVocabulary},
	}
	registry := registryFor(t, config)

	_, err := effective.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.vocabulary")
}

func TestResolvePacksRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{builtin.PackVocabulary},
		PackConfigs: policy.PackConfigs{
			builtin.PackVocabulary: vocabulary.EncodeConfig(vocabulary.Config{
				Go: vocabulary.GoConfig{ForbiddenTypeSuffixes: []string{"Repository"}},
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
		EnabledPacks: []string{builtin.PackText},
		PackConfigs: policy.PackConfigs{
			builtin.PackText: text.EncodeConfig(text.Config{}),
		},
	}
	registry := registryFor(t, config)

	_, err := effective.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.text.section_headers")
}

func TestResolvePacksRejectsUnsupportedConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{builtin.PackMarkdown},
		PackConfigs: policy.PackConfigs{
			builtin.PackMarkdown: {"unknown": true},
		},
	}
	registry := registryFor(t, config)

	_, err := effective.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.markdown config is not supported")
}

func TestResolvePacksAcceptsPacksWithoutConfig(t *testing.T) {
	t.Parallel()

	config := policy.Config{
		EnabledPacks: []string{builtin.PackMarkdown},
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
		EnabledPacks: []string{builtin.PackBash},
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
		EnabledPacks: []string{builtin.PackBash},
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

	registry, err := builtin.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	return registry
}
