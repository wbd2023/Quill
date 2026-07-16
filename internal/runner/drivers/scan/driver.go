package scan

import (
	"errors"
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func repositoryScanDriver(scanners runtimebinding.RepositoryScanners) (driver runner.Driver) {
	return func(
		context runner.Context,
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

		return scanner(context, execution)
	}
}
