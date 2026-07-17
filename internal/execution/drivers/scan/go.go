package scan

import (
	"context"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
)

// CheckGoArchitecture check go architecture.
func CheckGoArchitecture(goPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanGoArchitecture(ctx, context, execution, goPackID)
	}
}
