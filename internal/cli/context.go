package cli

import (
	"path/filepath"

	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
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

	compiled, err := profile.Compile(config, registry)
	if err != nil {
		return runner.Context{}, err
	}

	layout := runtime.NewLayout(repoRoot)
	goEnvironment := layout.GoEnvironment()
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDirectory(), "golangci")

	return runner.NewContext(
		repoRoot,
		scope,
		compiled.Profile,
		compiled.Effective,
		registry.ToolCapabilities(),
		layout.ToolEnvironment(),
		goEnvironment,
	), nil
}
