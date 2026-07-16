package cli

import (
	"path/filepath"

	"ciphera/tools/internal/ecosystem/golang"
	"ciphera/tools/internal/ecosystem/node"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/workspace"
)

func loadContext(repoRoot string, scope style.Scope) (context runner.Context, err error) {
	config, err := profile.Load(repoRoot)
	if err != nil {
		return runner.Context{}, err
	}

	if scope == "" {
		scope = config.Repository.DefaultScope
	}

	if !config.Repository.HasScope(scope) {
		return runner.Context{}, errUnknownScope(scope)
	}

	registry, err := shipped.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		return runner.Context{}, err
	}

	config, err = pack.ResolvePacks(config, registry.Packs())
	if err != nil {
		return runner.Context{}, err
	}

	compiled, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		return runner.Context{}, err
	}

	layout := workspace.NewLayout(repoRoot)
	path := layout.BuildPath(node.BinaryDirectory(layout))
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
	), nil
}
