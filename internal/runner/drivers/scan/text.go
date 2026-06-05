package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/runner"
)

func textPackScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerASCII: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return text.CheckASCII(context.RepoRoot, context.Profile.Repository, context.Scope)
		},
		builtin.ScannerExceptionMarkers: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return text.CheckExceptionMarkers(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
		builtin.ScannerLineLength: scanLineLengths,
		builtin.ScannerMaintenanceMarkers: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return text.CheckMaintenanceMarkers(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
		builtin.ScannerSectionHeaderNames:   scanSectionHeaderNames,
		builtin.ScannerSectionHeaderDensity: scanSectionHeaderDensity,
		builtin.ScannerSectionHeaders:       scanSectionHeaders,
	}
}

func scanLineLengths(
	context runner.Context,
	execution contract.RepositoryScanExecution,
) (result contract.ExecutionResult, err error) {
	files, err := runner.CollectFileSetFiles(context, execution.FileSet)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	return text.CheckLineLengths(context.RepoRoot, files)
}
