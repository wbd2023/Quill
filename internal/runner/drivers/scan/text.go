package scan

import (
	"ciphera/tools/internal/checks/text"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
)

func CheckASCII() (scanner binding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return text.CheckASCII(context.RepoRoot, context.Profile.Repository, context.Scope)
	}
}

func CheckExceptionMarkers() (scanner binding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return text.CheckExceptionMarkers(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}

func CheckLineLengths() (scanner binding.RepositoryScanner) {
	return scanLineLengths
}

func CheckMaintenanceMarkers() (scanner binding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return text.CheckMaintenanceMarkers(
			context.RepoRoot,
			context.Profile.Repository,
			context.Scope,
		)
	}
}

func CheckSectionHeaderNames(textPackID string) (scanner binding.RepositoryScanner) {
	return func(
		context runner.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderNames(context, execution, textPackID)
	}
}

func CheckSectionHeaderDensity(textPackID string) (scanner binding.RepositoryScanner) {
	return func(
		context runner.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderDensity(context, execution, textPackID)
	}
}

func CheckSectionHeaders(textPackID string) (scanner binding.RepositoryScanner) {
	return func(
		context runner.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaders(context, execution, textPackID)
	}
}

func scanLineLengths(
	context runner.Context,
	execution style.RepositoryScanExecution,
) (result style.ExecutionResult, err error) {
	files, err := runner.CollectFileSetFiles(context, execution.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return text.CheckLineLengths(context.RepoRoot, files)
}
