package execution

import (
	"context"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

func TestRunRuleUsesInjectedDriver(t *testing.T) {
	repoRoot := t.TempDir()
	rule := style.Rule{
		ID: "test/rule",
		Check: style.RepositoryScan{
			Scanner: "test",
		},
	}
	runCtx := NewRunContext(
		repoRoot,
		style.Scope("all"),
		profile.Profile{},
		style.Plan{},
		nil,
		nil,
		nil,
	)
	drivers := DriverSet{
		RepositoryScan: func(
			_ context.Context,
			_ RunContext,
			_ style.Rule,
			_ style.Job,
			_ toolchain.StatusMap,
		) (result style.ExecutionResult, err error) {
			return style.ExecutionResult{Diagnostics: []style.Diagnostic{{Message: "ran"}}}, nil
		},
	}

	result, err := RunRule(context.Background(), rule, runCtx, nil, drivers)
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
		Check: style.ToolchainCheck{
			ToolIDs: []string{"go"},
		},
	}
	runCtx := NewRunContext(
		repoRoot,
		style.Scope("all"),
		profile.Profile{},
		style.Plan{},
		nil,
		nil,
		nil,
	)
	drivers := DriverSet{}

	_, err := RunRule(context.Background(), rule, runCtx, nil, drivers)
	if err == nil {
		t.Fatal("expected error for execution with no registered driver, got nil")
	}
}
