package target

import (
	"fmt"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func targetCommandDriver(commands driverkit.TargetCommands) (driver execution.Executor) {
	return func(
		context execution.Context,
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

		return command(context, job)
	}
}

func targetCheckDriver(checks driverkit.TargetChecks) (driver execution.Executor) {
	return func(
		context execution.Context,
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

		return check(context, job)
	}
}
