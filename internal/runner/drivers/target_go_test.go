package drivers

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/pack/builtin"
)

func TestRunGolangciRulePassesCurrentAppScope(t *testing.T) {
	context := testContext(t, fixtures.RepositoryRoot(t), contract.Scope("app"))

	spec := contract.ExecutionSpec{
		Kind: contract.ExecutionTargetCommand,
		Detail: contract.TargetCommandExecution{
			ToolIDs:  []string{builtin.ToolGo, builtin.ToolGoimports, builtin.ToolGolangciLint},
			Action:   builtin.TargetActionGolangci,
			Language: builtin.LanguageGo,
			Targets:  []string{"app_go"},
		},
	}

	result, err := targetCommandDriver(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciDriver(app): %v\n%s", err, result.Output)
	}

	if result.Output != "0 issues." {
		t.Fatalf("unexpected app lint output: %q", result.Output)
	}
}

func TestRunGolangciRulePassesCurrentToolsScope(t *testing.T) {
	context := testContext(t, fixtures.RepositoryRoot(t), contract.Scope("tools"))

	spec := contract.ExecutionSpec{
		Kind: contract.ExecutionTargetCommand,
		Detail: contract.TargetCommandExecution{
			ToolIDs:  []string{builtin.ToolGo, builtin.ToolGoimports, builtin.ToolGolangciLint},
			Action:   builtin.TargetActionGolangci,
			Language: builtin.LanguageGo,
			Targets:  []string{"tools_go"},
		},
	}

	result, err := targetCommandDriver(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciDriver(tools): %v\n%s", err, result.Output)
	}

	if result.Output != "0 issues." {
		t.Fatalf("unexpected tools lint output: %q", result.Output)
	}
}
