package drivers

import (
	"context"
	"fmt"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

// profileCheckDriver returns the profile driver for check execution.
func profileCheckDriver(checks ProfileChecks) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		execution, found := job.(style.ProfileExecution)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf("profile driver received empty job")
		}

		check, found := checks.Lookup(execution.PackID, execution.Check)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown profile check %q for pack %q",
				execution.Check,
				execution.PackID,
			)
		}

		result, err = check(ctx, context, execution)
		result.PackID = execution.PackID
		return result, err
	}
}
