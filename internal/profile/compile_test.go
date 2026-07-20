package profile_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/pack/shipped"
	"github.com/wbd2023/Quill/internal/profile"
	"github.com/wbd2023/Quill/internal/testutil"
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

	config, err = pack.ResolvePacks(config, registry.Packs())
	if err != nil {
		t.Fatalf("ResolvePacks: %v", err)
	}

	compiled, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	if len(compiled.Effective.Rules) != len(config.Rules) {
		t.Fatalf(
			"expected %d effective rules, got %d",
			len(config.Rules),
			len(compiled.Effective.Rules),
		)
	}

	if _, found := compiled.Profile.FileSets.Lookup("line_length"); !found {
		t.Fatal("expected compiled profile to include Text Pack default file sets")
	}
}
