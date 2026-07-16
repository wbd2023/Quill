package scan

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// CheckGoArchitecture check go architecture.
func CheckGoArchitecture(goPackID string) (scanner runtimebinding.RepositoryScanner) {
	return func(
		context execution.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanGoArchitecture(context, execution, goPackID)
	}
}
