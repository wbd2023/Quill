package target

import (
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func targetCommandDriver(commands binding.TargetCommands) (driver runner.Driver) {
	return func(
		context runner.Context,
		spec style.ExecutionSpec,
		_ map[string]toolchain.Status,
	) (result style.ExecutionResult, err error) {
		execution, found := spec.TargetCommandExecution()
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"target command driver received empty spec",
			)
		}

		command, found := commands.Lookup(execution.Action)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown target command action %q",
				execution.Action,
			)
		}

		return command(context, spec)
	}
}

func targetCheckDriver(checks binding.TargetChecks) (driver runner.Driver) {
	return func(
		context runner.Context,
		spec style.ExecutionSpec,
		_ map[string]toolchain.Status,
	) (result style.ExecutionResult, err error) {
		execution, found := spec.TargetCheckExecution()
		if !found {
			return style.ExecutionResult{}, fmt.Errorf("target check driver received empty spec")
		}

		check, found := checks.Lookup(execution.Language)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unsupported target check language %q",
				execution.Language,
			)
		}

		return check(context, spec)
	}
}
