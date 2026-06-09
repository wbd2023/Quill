package runner

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func TestRunRuleUsesInjectedDriver(t *testing.T) {
	repoRoot := t.TempDir()
	rule := style.Rule{
		ID: "test/rule",
		Check: style.ExecutionSpec{
			Kind: style.ExecutionKind("test_execution"),
			Detail: style.RepositoryScanExecution{
				Scanner: "test",
			},
		},
	}
	context := NewContext(
		repoRoot,
		style.Scope("all"),
		policy.Config{},
		style.EffectiveConfig{},
		nil,
		nil,
		nil,
	)
	drivers := DriverRegistry{
		"test_execution": func(
			_ Context,
			_ style.ExecutionSpec,
			_ map[string]toolchain.Status,
		) (style.ExecutionResult, error) {
			return style.ExecutionResult{Output: "ran"}, nil
		},
	}

	result, err := RunRule(rule, context, nil, drivers)
	if err != nil {
		t.Fatalf("RunRule: %v", err)
	}

	if result.Output != "ran" {
		t.Fatalf("output = %q, want ran", result.Output)
	}
}
