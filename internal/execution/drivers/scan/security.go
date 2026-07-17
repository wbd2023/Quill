package scan

import (
	"context"

	"ciphera/tools/internal/checks/security"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
)

// CheckSecrets check secrets.
func CheckSecrets() (scanner driverkit.RepositoryScanner) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return security.CheckSecrets(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}
