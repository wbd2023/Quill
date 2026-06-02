package profile_test

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/profile"
)

func TestCompileResolvesCurrentProfileEnabledPacks(t *testing.T) {
	t.Parallel()

	config, err := profile.Load(fixtures.RepositoryRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := builtin.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	definitions := registry.Definitions()

	compiled, err := profile.Compile(config, registry)
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
