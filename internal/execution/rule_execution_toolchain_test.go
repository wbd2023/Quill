package execution

import (
	"context"
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

func TestRunRuleDispatchesToolchainJobWithInvalidTool(t *testing.T) {
	rule := style.Rule{
		ID: "toolchain/health",
		Check: style.ToolchainCheck{
			ToolIDs: []string{"go"},
		},
	}
	run := NewRunContext(
		t.TempDir(),
		style.Scope("all"),
		profile.Profile{},
		style.Plan{},
		nil,
		nil,
		nil,
	)
	statuses := toolchain.StatusMap{
		"go": {
			Tool:  toolchain.Tool{ID: "go", Name: "Go"},
			Valid: false,
			Issue: "missing",
		},
	}
	called := false
	drivers := DriverSet{
		Toolchain: func(
			_ context.Context,
			_ RunContext,
			_ style.Rule,
			_ style.Job,
			_ toolchain.StatusMap,
		) (style.ExecutionResult, error) {
			called = true
			return style.ExecutionResult{}, nil
		},
	}

	_, err := RunRule(context.Background(), rule, run, statuses, drivers)
	if err != nil {
		t.Fatalf("RunRule: %v", err)
	}
	if !called {
		t.Fatal("expected Toolchain driver to run")
	}
}
