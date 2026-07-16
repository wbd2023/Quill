package target

import (
	"testing"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
)

func TestRunGolangciRulePassesCurrentAppScope(t *testing.T) {
	context := testContext(t, testutil.RepositoryRoot(t), style.Scope("app"))

	job := style.TargetCommandJob{
		ToolIDs:  []string{tool.Go, tool.Goimports, tool.GolangciLint},
		Action:   golang.TargetActionGolangci,
		Language: golang.Language,
		Targets:  []string{"app_go"},
	}

	result, err := testTargetCommandDriver()(context, job, nil)
	if err != nil {
		t.Fatalf("golangciDriver(app): %v", err)
	}

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected app lint diagnostics: %#v", result.Diagnostics)
	}
}

func TestRunGolangciRulePassesCurrentToolsScope(t *testing.T) {
	context := testContext(t, testutil.RepositoryRoot(t), style.Scope("tools"))

	job := style.TargetCommandJob{
		ToolIDs:  []string{tool.Go, tool.Goimports, tool.GolangciLint},
		Action:   golang.TargetActionGolangci,
		Language: golang.Language,
		Targets:  []string{"tools_go"},
	}

	result, err := testTargetCommandDriver()(context, job, nil)
	if err != nil {
		t.Fatalf("golangciDriver(tools): %v", err)
	}

	if len(result.Diagnostics) != 0 {
		t.Fatalf("unexpected tools lint diagnostics: %#v", result.Diagnostics)
	}
}

func testTargetCommandDriver() (driver execution.Driver) {
	commands := runtimebinding.NewTargetCommands()
	commands.Add(
		golang.TargetActionGolangci,
		RunGolangci(golang.PackID, tool.GolangciLint, tool.Goimports, golang.Language),
	)
	return targetCommandDriver(commands)
}
