package execution

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
		Check: style.RepositoryScanExecution{
			Scanner: "test",
		},
	}
	context := NewRunContext(
		repoRoot,
		style.Scope("all"),
		policy.Config{},
		style.Plan{},
		nil,
		nil,
		nil,
	)
	drivers := ExecutorSet{
		RepositoryScan: func(
			_ RunContext,
			_ style.Job,
			_ toolchain.StatusMap,
		) (result style.ExecutionResult, err error) {
			return style.ExecutionResult{Diagnostics: []style.Diagnostic{{Message: "ran"}}}, nil
		},
	}

	result, err := RunRule(rule, context, nil, drivers)
	if err != nil {
		t.Fatalf("RunRule: %v", err)
	}

	if len(result.Diagnostics) == 0 || result.Diagnostics[0].Message != "ran" {
		t.Fatalf("diagnostics = %#v, want ran", result.Diagnostics)
	}
}

func TestRunRuleErrorsOnMissingDriver(t *testing.T) {
	repoRoot := t.TempDir()
	rule := style.Rule{
		ID: "test/unsupported",
		Check: style.ToolchainExecution{
			ToolIDs: []string{"go"},
		},
	}
	context := NewRunContext(
		repoRoot,
		style.Scope("all"),
		policy.Config{},
		style.Plan{},
		nil,
		nil,
		nil,
	)
	drivers := ExecutorSet{}

	_, err := RunRule(rule, context, nil, drivers)
	if err == nil {
		t.Fatal("expected error for execution with no registered driver, got nil")
	}
}
