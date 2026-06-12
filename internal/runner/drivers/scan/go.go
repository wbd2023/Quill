package scan

import (
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

func CheckGoArchitecture(goPackID string) (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanGoArchitecture(context, execution, goPackID)
	}
}
