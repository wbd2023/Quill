package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/rules/security"
	"ciphera/tools/internal/runner"
)

func securityRepositoryScanners() (scanners map[string]repositoryScanner) {
	return map[string]repositoryScanner{
		builtin.ScannerSecrets: func(
			context runner.Context,
			_ contract.ExecutionSpec,
		) (contract.ExecutionResult, error) {
			return security.CheckSecrets(
				context.RepoRoot,
				context.Profile.Repository,
				context.Scope,
			)
		},
	}
}
