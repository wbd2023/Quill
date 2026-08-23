package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/pack"
	"github.com/wbd2023/quill/internal/pack/shipped"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/testutil"
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

	if len(compiled.Rules) != len(config.Rules) {
		t.Fatalf(
			"expected %d compiled rules, got %d",
			len(config.Rules),
			len(compiled.Rules),
		)
	}

	if _, found := config.FileSets.Lookup("line_length"); !found {
		t.Fatal("expected resolved profile to include Text Pack default file sets")
	}
}
