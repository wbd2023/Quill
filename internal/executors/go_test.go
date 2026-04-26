package executors

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/rulepack"
)

func TestRunGolangciRulePassesCurrentAppScope(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.Scope("app"))

	spec := contract.ExecutionSpec{
		Kind: rulepack.ExecutorBackendCommand,
		Detail: contract.BackendCommandExecution{
			ToolIDs:  []string{rulepack.ToolGo, rulepack.ToolGoimports, rulepack.ToolGolangciLint},
			Action:   rulepack.BackendActionGolangci,
			Language: rulepack.LanguageGo,
			Backends: []string{"application_go"},
		},
	}

	result, err := backendCommandExecutor(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciExecutor(app): %v\n%s", err, result.Output)
	}

	if result.Output != "0 issues." {
		t.Fatalf("unexpected app lint output: %q", result.Output)
	}
}

func TestRunGolangciRulePassesCurrentToolsScope(t *testing.T) {
	context := testContext(t, fixtures.RepoRoot(t), contract.Scope("tools"))

	spec := contract.ExecutionSpec{
		Kind: rulepack.ExecutorBackendCommand,
		Detail: contract.BackendCommandExecution{
			ToolIDs:  []string{rulepack.ToolGo, rulepack.ToolGoimports, rulepack.ToolGolangciLint},
			Action:   rulepack.BackendActionGolangci,
			Language: rulepack.LanguageGo,
			Backends: []string{"tooling_go"},
		},
	}

	result, err := backendCommandExecutor(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciExecutor(tools): %v\n%s", err, result.Output)
	}

	if result.Output != "0 issues." {
		t.Fatalf("unexpected tools lint output: %q", result.Output)
	}
}
