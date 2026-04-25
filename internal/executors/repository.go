package executors

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	bashstyle "ciphera/tools/internal/rules/bash"
	gostyle "ciphera/tools/internal/rules/go"
	repostyle "ciphera/tools/internal/rules/repo"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

/* ------------------------------------- Repository Scanners ------------------------------------ */

func repositoryScanExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]runtime.ToolStatus,
) (output string, err error) {
	scanner, found := repositoryScanners()[spec.Scanner]
	if !found {
		return "", errors.New("unknown repository scanner")
	}

	return scanner(context, spec)
}

type repositoryScanner func(
	context runner.Context,
	spec contract.ExecutionSpec,
) (output string, err error)

func repositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		rulepack.RepositoryScannerArchitecture: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return gostyle.CheckArchitecture(context.RepoRoot, context.Scope, context.Policy)
		},
		rulepack.RepositoryScannerASCII: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return repostyle.CheckASCII(context.RepoRoot, context.Policy.Repository, context.Scope)
		},
		rulepack.RepositoryScannerBashMagicValues: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return bashstyle.CheckMagicValues(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerBashSafety: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return bashstyle.CheckSafety(context.RepoRoot, context.Policy.Repository, context.Scope)
		},
		rulepack.RepositoryScannerBashStructure: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return bashstyle.CheckStructure(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerBashTestHygiene: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return bashstyle.CheckTestHygiene(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerExceptionMarkers: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return repostyle.CheckExceptionMarkers(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerGuardClauseSpacing: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return gostyle.CheckGuardClauseSpacing(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerLineLength: func(
			context runner.Context,
			spec contract.ExecutionSpec,
		) (string, error) {
			files, err := runner.CollectFileSetFiles(context, spec.FileSet)
			if err != nil {
				return "", err
			}

			return repostyle.CheckLineLengths(context.RepoRoot, files)
		},
		rulepack.RepositoryScannerMaintenanceMarkers: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return repostyle.CheckMaintenanceMarkers(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerNaming: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return repostyle.CheckNaming(
				context.RepoRoot,
				context.Policy.Repository,
				context.Policy.Naming,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerSecrets: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return repostyle.CheckSecrets(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerSectionHeaderNames: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return repostyle.CheckSectionHeaderNames(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerSectionHeaders: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return repostyle.CheckSectionHeaders(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
		rulepack.RepositoryScannerSwitchCaseSpacing: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (string, error) {
			return gostyle.CheckSwitchCaseSpacing(
				context.RepoRoot,
				context.Policy.Repository,
				context.Scope,
			)
		},
	}
}
