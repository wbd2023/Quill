package runner

import (
	"testing"

	"ciphera/tools/internal/lockfile"
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
	context := NewContext(
		repoRoot,
		style.Scope("all"),
		policy.Config{},
		style.EffectiveConfig{},
		nil,
		nil,
		nil,
		lockfile.Lockfile{},
	)
	drivers := DriverSet{
		RepositoryScan: func(
			_ Context,
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
	context := NewContext(
		repoRoot,
		style.Scope("all"),
		policy.Config{},
		style.EffectiveConfig{},
		nil,
		nil,
		nil,
		lockfile.Lockfile{},
	)
	drivers := DriverSet{}

	_, err := RunRule(rule, context, nil, drivers)
	if err == nil {
		t.Fatal("expected error for execution with no registered driver, got nil")
	}
}
