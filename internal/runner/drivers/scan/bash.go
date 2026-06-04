package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/bash"
	"ciphera/tools/internal/runner"
)

func bashPackScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerBashMagicValues: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return bash.CheckMagicValues(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
		builtin.ScannerBashSafety: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return bash.CheckSafety(context.RepoRoot, context.Profile.Repository, context.Scope)
		},
		builtin.ScannerBashStructure: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return bash.CheckStructure(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
		builtin.ScannerBashTestHygiene: func(
			context runner.Context,
			_ contract.RepositoryScanExecution,
		) (contract.ExecutionResult, error) {
			return bash.CheckTestHygiene(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
	}
}
