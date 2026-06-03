package target

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

func testContext(
	t *testing.T,
	repoRoot string,
	scope contract.Scope,
) (context runner.Context) {
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

	layout := runtime.LayoutForRepository(repoRoot)
	goEnvironment := layout.GoEnvironment()
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDir, "golangci")

	return runner.NewContext(
		repoRoot,
		scope,
		compiled.Profile,
		compiled.Effective,
		registry.ToolCapabilities(),
		layout.ToolEnvironment(),
		goEnvironment,
	)
}
