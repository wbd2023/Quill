package executors

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

func goRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		rulepack.ScannerArchitecture: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runGoArchitectureCheck(context)
		},
	}
}
