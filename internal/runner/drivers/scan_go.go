package drivers

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/runner"
)

func goRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerArchitecture: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runGoArchitectureCheck(context)
		},
	}
}
