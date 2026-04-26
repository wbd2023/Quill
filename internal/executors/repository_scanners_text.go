package executors

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/runner"
)

func textRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		rulepack.ScannerASCII: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckASCII(context.RepoRoot, context.Policy.Repository, context.Scope)
		},
		rulepack.ScannerExceptionMarkers: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckExceptionMarkers(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.ScannerLineLength: func(
			context runner.Context,
			spec contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runLineLengthScanner(context, spec)
		},
		rulepack.ScannerMaintenanceMarkers: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckMaintenanceMarkers(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.ScannerSectionHeaderNames: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckSectionHeaderNames(
				context.RepoRoot,
				context.Policy.Repository,
				context.Policy.Formatting.SectionHeaders,
				context.Scope,
			)
		},
		rulepack.ScannerSectionHeaderDensity: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckSectionHeaderDensity(
				context.RepoRoot,
				context.Policy.Repository,
				context.Policy.Formatting.SectionHeaders,
				context.Scope,
			)
		},
		rulepack.ScannerSectionHeaders: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckSectionHeaders(
				context.RepoRoot,
				context.Policy.Repository,
				context.Policy.Formatting.SectionHeaders,
				context.Scope,
			)
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
