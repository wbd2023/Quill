package executors

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/rules/security"
	"ciphera/tools/internal/runner"
)

func securityRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
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
	}
}
