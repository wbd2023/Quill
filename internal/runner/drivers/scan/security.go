package scan

import (
	"ciphera/tools/internal/checks/security"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

func CheckSecrets() (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return security.CheckSecrets(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}
