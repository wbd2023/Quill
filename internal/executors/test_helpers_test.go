package executors

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

func testContext(
	t *testing.T,
	repoRoot string,
	scope contract.Scope,
) (context runner.Context) {
	t.Helper()

	policy, err := profile.Load(repoRoot)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	registry, err := rulepack.DefaultRegistry(policy.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	effective, err := policy.Compile(registry)
	if err != nil {
		t.Fatalf("Compile: %v", err)
	}

	return runner.NewContext(repoRoot, scope, policy, effective)
}
