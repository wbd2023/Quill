package scan

import (
	"context"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
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
