package cli

import (
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
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

	if !config.Repository.ScopeExists(scope) {
		return runner.Context{}, errUnknownScope(scope)
	}

	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		return runner.Context{}, err
	}

	effective, err := profile.Compile(config, registry.Definitions())
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
		effective,
		registry.ToolCapabilities(),
		layout.ToolEnvironment(),
		goEnvironment,
	), nil
}
