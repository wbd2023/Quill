package drivers

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/runner"
)

func textRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerASCII: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckASCII(context.RepoRoot, context.Profile.Repository, context.Scope)
		},
		builtin.ScannerExceptionMarkers: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckExceptionMarkers(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
		builtin.ScannerLineLength: func(
			context runner.Context,
			spec contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runLineLengthScanner(context, spec)
		},
		builtin.ScannerMaintenanceMarkers: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckMaintenanceMarkers(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
		builtin.ScannerSectionHeaderNames: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runSectionHeaderNamesScanner(context)
		},
		builtin.ScannerSectionHeaderDensity: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runSectionHeaderDensityScanner(context)
		},
		builtin.ScannerSectionHeaders: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runSectionHeadersScanner(context)
		},
	}
}

func runLineLengthScanner(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.RepositoryScanExecution()
	if !found {
		return contract.ExecutionResult{}, errors.New(
			"line-length scanner received empty spec",
		)
	}

	files, err := runner.CollectFileSetFiles(context, execution.FileSet)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	return text.CheckLineLengths(context.RepoRoot, files)
}
