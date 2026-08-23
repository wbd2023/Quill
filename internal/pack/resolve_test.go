package pack_test

import (
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/pack/shipped/bash"
	"github.com/wbd2023/quill/internal/pack/shipped/markdown"
	"github.com/wbd2023/quill/internal/pack/shipped/text"
	textpolicy "github.com/wbd2023/quill/internal/pack/shipped/text/policy"
	"github.com/wbd2023/quill/internal/pack/shipped/vocabulary"
	vocabularypolicy "github.com/wbd2023/quill/internal/pack/shipped/vocabulary/policy"
	"github.com/wbd2023/quill/internal/profile"
)

/* ----------------------------------------- Validation ----------------------------------------- */

func TestResolvePacksRejectsMissingRequiredConfig(t *testing.T) {
	t.Parallel()

	config := profile.Profile{
		EnabledPacks: []string{vocabulary.PackID},
	}
	registry := registryFor(t, config)

	_, err := pack.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.vocabulary")
}

func TestResolvePacksRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	config := profile.Profile{
		EnabledPacks: []string{vocabulary.PackID},
		PackPolicies: profile.PackPolicies{
			vocabulary.PackID: vocabularypolicy.EncodeConfig(vocabularypolicy.Config{
				Go: vocabularypolicy.GoConfig{
					TypeSuffixes: map[string][]string{"": {"Repository"}},
				},
			}),
		},
	}
	registry := registryFor(t, config)

	_, err := pack.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.vocabulary.go.type_suffixes")
}

func TestResolvePacksRejectsInvalidTextConfig(t *testing.T) {
	t.Parallel()

	config := profile.Profile{
		EnabledPacks: []string{text.PackID},
		PackPolicies: profile.PackPolicies{
			text.PackID: textpolicy.EncodeConfig(textpolicy.Config{}),
		},
	}
	registry := registryFor(t, config)

	_, err := pack.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.text.section_headers")
}

func TestResolvePacksRejectsUnsupportedConfig(t *testing.T) {
	t.Parallel()

	config := profile.Profile{
		EnabledPacks: []string{markdown.PackID},
		PackPolicies: profile.PackPolicies{
			markdown.PackID: {"unknown": true},
		},
	}
	registry := registryFor(t, config)

	_, err := pack.ResolvePacks(config, registry.Packs())
	requireErrorContains(t, err, "packs.markdown policy is not supported")
}

func TestResolvePacksAcceptsPacksWithoutConfig(t *testing.T) {
	t.Parallel()

	config := profile.Profile{
		EnabledPacks: []string{markdown.PackID},
	}
	registry := registryFor(t, config)

	if _, err := pack.ResolvePacks(config, registry.Packs()); err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}
}

/* -------------------------------------- Default File Sets ------------------------------------- */

func TestResolvePacksAppliesPackDefaultFileSets(t *testing.T) {
	t.Parallel()

	config := profile.Profile{
		EnabledPacks: []string{bash.PackID},
	}
	registry := registryFor(t, config)

	resolved, err := pack.ResolvePacks(config, registry.Packs())
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

	config := profile.Profile{
		EnabledPacks: []string{bash.PackID},
		FileSets: profile.FileSets{
			{Name: "bash", Include: profile.FileSetInclude{Extensions: []string{".bash"}}},
		},
	}
	registry := registryFor(t, config)

	resolved, err := pack.ResolvePacks(config, registry.Packs())
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

func registryFor(t *testing.T, config profile.Profile) (registry pack.Registry) {
	t.Helper()

	registry, err := shipped.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	return registry
}

func requireErrorContains(tb testing.TB, err error, text string) {
	tb.Helper()

	if err == nil {
		tb.Fatalf("expected error containing %q, got nil", text)
	}

	if !strings.Contains(err.Error(), text) {
		tb.Fatalf("expected error containing %q, got %q", text, err.Error())
	}
}
