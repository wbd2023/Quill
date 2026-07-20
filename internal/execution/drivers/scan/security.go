package scan

import (
	"context"

	"github.com/wbd2023/Quill/internal/checks/security"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
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
