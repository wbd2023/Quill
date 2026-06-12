package target

import (
	"testing"

	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
)

func TestRunGolangciRulePassesCurrentAppScope(t *testing.T) {
	context := testContext(t, testutil.RepositoryRoot(t), style.Scope("app"))

	spec := style.ExecutionSpec{
		Kind: style.ExecutionTargetCommand,
		Detail: style.TargetCommandExecution{
			ToolIDs:  []string{tool.Go, tool.Goimports, tool.GolangciLint},
			Action:   golang.TargetActionGolangci,
			Language: golang.Language,
			Targets:  []string{"app_go"},
		},
	}

	result, err := testTargetCommandDriver()(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciDriver(app): %v\n%s", err, result.Output)
	}

	if result.Output != "0 issues." {
		t.Fatalf("unexpected app lint output: %q", result.Output)
	}
}

func TestRunGolangciRulePassesCurrentToolsScope(t *testing.T) {
	context := testContext(t, testutil.RepositoryRoot(t), style.Scope("tools"))

	spec := style.ExecutionSpec{
		Kind: style.ExecutionTargetCommand,
		Detail: style.TargetCommandExecution{
			ToolIDs:  []string{tool.Go, tool.Goimports, tool.GolangciLint},
			Action:   golang.TargetActionGolangci,
			Language: golang.Language,
			Targets:  []string{"tools_go"},
		},
	}

	result, err := testTargetCommandDriver()(context, spec, nil)
	if err != nil {
		t.Fatalf("golangciDriver(tools): %v\n%s", err, result.Output)
	}

	if result.Output != "0 issues." {
		t.Fatalf("unexpected tools lint output: %q", result.Output)
	}
}

func testTargetCommandDriver() (driver runner.Driver) {
	commands := runtimebinding.NewTargetCommands()
	commands.Add(
		golang.TargetActionGolangci,
		RunGolangci(golang.PackID, tool.GolangciLint, tool.Goimports, golang.Language),
	)
	return targetCommandDriver(commands)
}
