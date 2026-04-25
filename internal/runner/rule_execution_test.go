package runner

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/runtime"
)

func TestRunRuleUsesInjectedExecutor(t *testing.T) {
	rule := contract.Rule{
		RuleDefinition: contract.RuleDefinition{
			ID: "test/rule",
			Spec: contract.ExecutionSpec{
				Executor: "test_executor",
			},
		},
	}
	context := NewContext(
		t.TempDir(),
		contract.ScopeAll,
		profile.Profile{},
		profile.EffectiveConfig{},
	)
	executors := ExecutorRegistry{
		"test_executor": func(
			_ Context,
			_ contract.ExecutionSpec,
			_ map[string]runtime.ToolStatus,
		) (string, error) {
			return "ran", nil
		},
	}

	output, err := RunRule(rule, context, nil, executors)
	if err != nil {
		t.Fatalf("RunRule: %v", err)
	}

	if output != "ran" {
		t.Fatalf("output = %q, want ran", output)
	}
}
