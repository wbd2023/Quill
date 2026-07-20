package scan

import (
	"path/filepath"
	"testing"

	"github.com/wbd2023/Quill/internal/ecosystem/golang"
	"github.com/wbd2023/Quill/internal/ecosystem/node"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/pack/shipped"
	"github.com/wbd2023/Quill/internal/profile"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/workspace"
)

func testContext(
	t *testing.T,
	repoRoot string,
	scope style.Scope,
) (context execution.RunContext) {
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

	layout := workspace.NewLayout(repoRoot)
	path := layout.BuildPath(node.BinaryDirectory(layout))
	toolEnvironment := map[string]string{"PATH": path}
	goEnvironment := golang.Environment(layout, path)
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDirectory(), "golangci")

	return execution.NewRunContext(
		repoRoot,
		scope,
		compiled.Profile,
		compiled.Effective,
		registry.ToolCapabilities(),
		toolEnvironment,
		goEnvironment,
	)
}
