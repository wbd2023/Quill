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
		rule style.Rule,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		check, ok := job.(style.ProfileCheck)
		if !ok {
			return style.ExecutionResult{}, fmt.Errorf(
				"profile driver received unsupported execution job %T",
				job,
			)
		}

		bound, found := checks.Lookup(rule.PackID, check.Check)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown profile check %q for pack %q",
				check.Check,
				rule.PackID,
			)
		}

		return bound(ctx, context, check)
	}
}
