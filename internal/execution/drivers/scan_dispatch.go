package drivers

import (
	"context"
	"errors"
	"fmt"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

// repositoryScanDriver returns the repository-scan driver for check execution.
func repositoryScanDriver(scanners RepositoryScanners) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		execution, found := job.(style.RepositoryScanExecution)
		if !found {
			return style.ExecutionResult{},
				errors.New("repository-scan driver received empty job")
		}

		scanner, found := scanners.Lookup(execution.PackID, execution.Scanner)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown repository scanner %q for pack %q",
				execution.Scanner,
				execution.PackID,
			)
		}

		result, err = scanner(ctx, context, execution)
		result.PackID = execution.PackID
		return result, err
	}
}
