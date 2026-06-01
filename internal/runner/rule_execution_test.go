package runner

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/toolchain"
)

func TestRunRuleUsesInjectedExecutor(t *testing.T) {
	repoRoot := t.TempDir()
	rule := contract.Rule{
		ID: "test/rule",
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorKind("test_executor"),
			Detail: contract.RepositoryScanExecution{
				Scanner: "test",
			},
		},
	}
	context := NewContext(
		repoRoot,
		contract.Scope("all"),
		policy.Config{},
		contract.EffectiveConfig{},
		nil,
		nil,
		nil,
	)
	executors := ExecutorRegistry{
		"test_executor": func(
			_ Context,
			_ contract.ExecutionSpec,
			_ map[string]toolchain.Status,
		) (contract.ExecutionResult, error) {
			return contract.ExecutionResult{Output: "ran"}, nil
		},
	}

	result, err := RunRule(rule, context, nil, executors)
	if err != nil {
		t.Fatalf("RunRule: %v", err)
	}

	if result.Output != "ran" {
		t.Fatalf("output = %q, want ran", result.Output)
	}
}
