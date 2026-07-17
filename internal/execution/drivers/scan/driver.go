package scan

import (
	"context"
	"errors"
	"fmt"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func repositoryScanDriver(scanners driverkit.RepositoryScanners) (driver execution.Executor) {
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
