package executors

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/rules/bash"
	"ciphera/tools/internal/runner"
)

func bashRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
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
	}
}
