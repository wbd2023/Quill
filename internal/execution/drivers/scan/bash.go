package scan

import (
	"context"

	"github.com/wbd2023/Quill/internal/checks/bash"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
)

// CheckBashMagicValues check bash magic values.
func CheckBashMagicValues() (scanner driverkit.RepositoryScanner) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckMagicValues(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}

// CheckBashSafety check bash safety.
func CheckBashSafety() (scanner driverkit.RepositoryScanner) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckSafety(context.RepoRoot, context.Profile.Repository, context.Scope)
	}
}

// CheckBashStructure check bash structure.
func CheckBashStructure() (scanner driverkit.RepositoryScanner) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckStructure(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}

// CheckBashTestHygiene check bash test hygiene.
func CheckBashTestHygiene() (scanner driverkit.RepositoryScanner) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return bash.CheckTestHygiene(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}
