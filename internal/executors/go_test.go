package executors

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
)

func TestRunGolangciRulePassesCurrentAppScope(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.ScopeApp)

	spec := contract.ExecutionSpec{
		Backend: "go_app",
	}

	output, err := golangciExecutor(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciExecutor(app): %v\n%s", err, output)
	}

	if output != "0 issues.\n" {
		t.Fatalf("unexpected app lint output: %q", output)
	}
}

func TestRunGolangciRulePassesCurrentToolsScope(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.ScopeTools)

	spec := contract.ExecutionSpec{
		Backend: "go_tools",
	}

	output, err := golangciExecutor(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciExecutor(tools): %v\n%s", err, output)
	}

	if output != "0 issues.\n" {
		t.Fatalf("unexpected tools lint output: %q", output)
	}
}
