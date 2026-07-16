package scan

import (
	"path/filepath"
	"testing"

	"ciphera/tools/internal/ecosystem/golang"
	"ciphera/tools/internal/ecosystem/node"
	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
)

func testContext(
	t *testing.T,
	repoRoot string,
	scope style.Scope,
) (context runner.Context) {
	t.Helper()

	config, err := profile.Load(repoRoot)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := shipped.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	compiled, err := profile.Compile(config, registry)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	layout := runtime.NewLayout(repoRoot)
	path := runtime.SearchPath(layout.ToolBinaryDirectory(), node.BinaryDirectory(layout))
	toolEnvironment := map[string]string{"PATH": path}
	goEnvironment := golang.Environment(layout, path)
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDirectory(), "golangci")

	return runner.NewContext(
		repoRoot,
		scope,
		compiled.Profile,
		compiled.Effective,
		registry.ToolCapabilities(),
		toolEnvironment,
		goEnvironment,
		lockfile.Lockfile{},
	)
}
