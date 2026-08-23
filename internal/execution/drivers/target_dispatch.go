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
		rule style.Rule,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		command, ok := job.(style.TargetCommandJob)
		if !ok {
			return style.ExecutionResult{}, fmt.Errorf(
				"target command driver received unsupported execution job %T",
				job,
			)
		}

		bound, found := commands.Lookup(rule.PackID, command.Language, command.Action)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown target command %q for language %q in pack %q",
				command.Action,
				command.Language,
				rule.PackID,
			)
		}

		return bound(ctx, context, command)
	}
}

func targetCheckDriver(checks TargetChecks) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		rule style.Rule,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		check, ok := job.(style.TargetCheckJob)
		if !ok {
			return style.ExecutionResult{}, fmt.Errorf(
				"target check driver received unsupported execution job %T",
				job,
			)
		}

		bound, found := checks.Lookup(rule.PackID, check.Language, check.Check)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown target check %q for language %q in pack %q",
				check.Check,
				check.Language,
				rule.PackID,
			)
		}

		return bound(ctx, context, check)
	}
}
