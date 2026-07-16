package profile_test

import (
	"testing"

	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/testutil"
)

func TestCompileResolvesCurrentProfileEnabledPacks(t *testing.T) {
	t.Parallel()

	config, err := profile.Load(testutil.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := shipped.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	definitions := registry.Definitions()

	config, err = pack.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}

	compiled, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	if len(compiled.Effective.Rules) != len(definitions.Rules) {
		t.Fatalf(
			"expected %d effective rules, got %d",
			len(definitions.Rules),
			len(compiled.Effective.Rules),
		)
	}

	if _, found := compiled.Profile.FileSets.Lookup("bash"); !found {
		t.Fatal("expected compiled profile to include Pack default file sets")
	}
}
