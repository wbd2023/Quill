package runner

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
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

	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	return NewContext(
		repoRoot,
		scope,
		config,
		effective,
		registry.ToolCapabilities(),
		map[string]string{"PATH": ""},
		map[string]string{"PATH": ""},
	)
}
