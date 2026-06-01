package cli

import (
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/profile/effective"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

func loadContext(repoRoot string, scope contract.Scope) (context runner.Context, err error) {
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

	registry, err := builtin.DefaultRegistry(config.EnabledPacks)
	if err != nil {
		return runner.Context{}, err
	}

	config, err = effective.ResolvePacks(config, registry.Packs())
	if err != nil {
		return runner.Context{}, err
	}

	compiled, err := profile.Compile(config, registry.Definitions())
	if err != nil {
		return runner.Context{}, err
	}

	layout := runtime.LayoutForRepository(repoRoot)
	goEnvironment := layout.GoEnvironment()
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDir, "golangci")

	return runner.NewContext(
		repoRoot,
		scope,
		config,
		compiled,
		registry.ToolCapabilities(),
		layout.ToolEnvironment(),
		goEnvironment,
	), nil
}
