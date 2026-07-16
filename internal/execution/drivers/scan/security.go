package scan

import (
	"ciphera/tools/internal/checks/security"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// CheckSecrets check secrets.
func CheckSecrets() (scanner runtimebinding.RepositoryScanner) {
	return func(
		context execution.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return security.CheckSecrets(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}
