package executors

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/toolchain"
)

func repositoryScanExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.RepositoryScanExecution()
	if !found {
		return contract.ExecutionResult{},
			errors.New("repository-scan executor received empty spec")
	}

	scanner, found := repositoryScanners()[execution.Scanner]
	if !found {
		return contract.ExecutionResult{}, errors.New("unknown repository scanner")
	}

	return scanner(context, spec)
}

type repositoryScanner func(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error)
