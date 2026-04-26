package executors

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
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

	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	layout := runtime.LayoutForRepository(repoRoot)
	goEnvironment := layout.GoEnvironment()
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDir, "golangci")

	return runner.NewContext(
		repoRoot,
		scope,
		config,
		effective,
		registry.ToolCapabilities(),
		layout.ToolEnvironment(),
		goEnvironment,
	)
}
