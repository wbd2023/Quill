package execution

import (
	"testing"

	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/style"
)

func testContext(
	t *testing.T,
	repoRoot string,
	scope style.Scope,
) (context Context) {
	t.Helper()

	config, err := profile.Load(repoRoot)
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

	return NewContext(
		repoRoot,
		scope,
		compiled.Profile,
		compiled.Effective,
		registry.ToolCapabilities(),
		map[string]string{"PATH": ""},
		map[string]string{"PATH": ""},
	)
}
