package runner

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/profile"
)

func testContext(
	t *testing.T,
	repoRoot string,
	scope contract.Scope,
) (context Context) {
	t.Helper()

	config, err := profile.Load(repoRoot)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := builtin.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	compiled, err := profile.Compile(config, registry)
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
