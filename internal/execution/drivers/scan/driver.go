package scan

import (
	"context"
	"errors"
	"fmt"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

// CheckDriver returns the repository-scan driver for check execution.
func CheckDriver(scanners driverkit.RepositoryScanners) (driver execution.Driver) {
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

		scanner, found := scanners.Lookup(execution.Scanner)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown repository scanner %q",
				execution.Scanner,
			)
		}

		return scanner(ctx, context, execution)
	}
}
