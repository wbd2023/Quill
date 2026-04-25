package cli

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

func loadContext(repoRoot string, scope contract.Scope) (context runner.Context, err error) {
	policy, err := profile.Load(repoRoot)
	if err != nil {
		return runner.Context{}, err
	}

	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		return runner.Context{}, err
	}

	effective, err := policy.Compile(registry)
	if err != nil {
		return runner.Context{}, err
	}

	return runner.NewContext(repoRoot, scope, policy, effective), nil
}
