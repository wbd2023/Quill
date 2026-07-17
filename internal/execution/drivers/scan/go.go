package scan

import (
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
)

// CheckGoArchitecture check go architecture.
func CheckGoArchitecture(goPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.RunContext,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanGoArchitecture(context, execution, goPackID)
	}
}
