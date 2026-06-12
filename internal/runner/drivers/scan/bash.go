package scan

import (
	"ciphera/tools/internal/checks/bash"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

func CheckBashMagicValues() (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckMagicValues(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}

func CheckBashSafety() (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckSafety(context.RepoRoot, context.Profile.Repository, context.Scope)
	}
}

func CheckBashStructure() (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckStructure(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}

func CheckBashTestHygiene() (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckTestHygiene(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}
