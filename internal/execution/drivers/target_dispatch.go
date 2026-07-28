package drivers

import (
	"context"
	"fmt"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

func targetCommandDriver(commands TargetCommands) (driver execution.Driver) {
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

		command, found := commands.Lookup(execution.PackID, execution.Language, execution.Action)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown target command %q for language %q in pack %q",
				execution.Action,
				execution.Language,
				execution.PackID,
			)
		}

		result, err = command(ctx, context, execution)
		result.PackID = execution.PackID
		return result, err
	}
}

func targetCheckDriver(checks TargetChecks) (driver execution.Driver) {
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

		check, found := checks.Lookup(execution.PackID, execution.Language, execution.Check)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown target check %q for language %q in pack %q",
				execution.Check,
				execution.Language,
				execution.PackID,
			)
		}

		result, err = check(ctx, context, execution)
		result.PackID = execution.PackID
		return result, err
	}
}
