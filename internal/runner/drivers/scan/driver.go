package scan

import (
	"errors"
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func repositoryScanDriver(scanners binding.RepositoryScanners) (driver runner.Driver) {
	return func(
		context runner.Context,
		spec style.ExecutionSpec,
		_ map[string]toolchain.Status,
	) (result style.ExecutionResult, err error) {
		execution, found := spec.RepositoryScanExecution()
		if !found {
			return style.ExecutionResult{},
				errors.New("repository-scan driver received empty spec")
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
