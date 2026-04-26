package executors

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/rules/bash"
	"ciphera/tools/internal/rules/naming"
	"ciphera/tools/internal/rules/security"
	"ciphera/tools/internal/rules/text"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

func repositoryScanExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	detail, found := spec.RepositoryScanExecution()
	if !found {
		return contract.ExecutionResult{},
			errors.New("repository-scan executor received empty spec")
	}

	scanner, found := repositoryScanners()[detail.Scanner]
	if !found {
		return contract.ExecutionResult{}, errors.New("unknown repository scanner")
	}

	return scanner(context, spec)
}

type repositoryScanner func(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error)

func repositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		rulepack.ScannerArchitecture: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return runGoArchitectureCheck(context)
		},
		rulepack.ScannerASCII: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return text.CheckASCII(context.RepoRoot, context.Policy.Repository, context.Scope)
		},
		rulepack.ScannerBashMagicValues: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return bash.CheckMagicValues(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.ScannerBashSafety: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return bash.CheckSafety(context.RepoRoot, context.Policy.Repository, context.Scope)
		},
		rulepack.ScannerBashStructure: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return bash.CheckStructure(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.ScannerBashTestHygiene: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return bash.CheckTestHygiene(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
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
			detail, found := spec.RepositoryScanExecution()
			if !found {
				return contract.ExecutionResult{}, errors.New(
					"line-length scanner received empty spec",
				)
			}

			files, err := runner.CollectFileSetFiles(context, detail.FileSet)
			if err != nil {
				return contract.ExecutionResult{}, err
			}

			return text.CheckLineLengths(context.RepoRoot, files)
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
		rulepack.ScannerNaming: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return naming.CheckNaming(
				context.RepoRoot,
				context.Policy.Repository,
				context.Policy.Naming,
				context.Scope,
			)
		},
		rulepack.ScannerSecrets: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return security.CheckSecrets(
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
