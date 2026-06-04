package scan

import (
	"errors"
	"fmt"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/toolchain"
)

func repositoryScanDriver(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.RepositoryScanExecution()
	if !found {
		return contract.ExecutionResult{},
			errors.New("repository-scan driver received empty spec")
	}

	scanner, found := scanners[execution.Scanner]
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf(
			"unknown repository scanner %q",
			execution.Scanner,
		)
	}

	return scanner(context, execution)
}

type repositoryScanner func(
	context runner.Context,
	execution contract.RepositoryScanExecution,
) (result contract.ExecutionResult, err error)
