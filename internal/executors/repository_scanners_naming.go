package executors

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/rules/naming"
	"ciphera/tools/internal/runner"
)

func namingRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		rulepack.ScannerNaming: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return naming.CheckNaming(
				context.RepoRoot,
				context.Policy.Repository,
				context.Policy.Vocabulary,
				context.Scope,
			)
		},
	}
}
