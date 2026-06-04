package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/runner"
)

func goPackScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerArchitecture: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return scanGoArchitecture(context)
		},
	}
}
