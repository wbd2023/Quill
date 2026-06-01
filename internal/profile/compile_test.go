package profile_test

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/profile/effective"
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

	config, err = effective.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("effective.ResolvePacks: %v", err)
	}

	definitions := registry.Definitions()

	compiled, err := profile.Compile(config, definitions)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	if len(compiled.Rules) != len(definitions.Rules) {
		t.Fatalf(
			"expected %d effective rules, got %d",
			len(definitions.Rules),
			len(compiled.Rules),
		)
	}
}
