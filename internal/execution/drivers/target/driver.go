package target

import (
	"context"
	"fmt"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

func targetCommandDriver(commands driverkit.TargetCommands) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		execution, found := job.(style.TargetCommandJob)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"target command driver received empty job",
			)
		}

		command, found := commands.Lookup(execution.Action)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown target command action %q",
				execution.Action,
			)
		}

		return command(ctx, context, job)
	}
}

func targetCheckDriver(checks driverkit.TargetChecks) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		execution, found := job.(style.TargetCheckJob)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf("target check driver received empty job")
		}

		check, found := checks.Lookup(execution.Language)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unsupported target check language %q",
				execution.Language,
			)
		}

		return check(ctx, context, job)
	}
}

// CheckCommandDriver returns the target-command driver for check execution.
func CheckCommandDriver(commands driverkit.TargetCommands) (driver execution.Driver) {
	return targetCommandDriver(commands)
}

// CheckCheckDriver returns the target-check driver for check execution.
func CheckCheckDriver(checks driverkit.TargetChecks) (driver execution.Driver) {
	return targetCheckDriver(checks)
}

// FixCommandDriver returns the target-command driver for fix execution.
func FixCommandDriver(commands driverkit.TargetCommands) (driver execution.Driver) {
	return targetCommandDriver(commands)
}
