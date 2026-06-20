package scan

import (
	"ciphera/tools/internal/checks/text"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// CheckASCII check a s c i i.
func CheckASCII() (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return text.CheckASCII(context.RepoRoot, context.Profile.Repository, context.Scope)
	}
}

// CheckExceptionMarkers check exception markers.
func CheckExceptionMarkers() (scanner runtimebinding.RepositoryScanner) {
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

// CheckLineLengths check line lengths.
func CheckLineLengths() (scanner runtimebinding.RepositoryScanner) {
	return scanLineLengths
}

// CheckMaintenanceMarkers check maintenance markers.
func CheckMaintenanceMarkers() (scanner runtimebinding.RepositoryScanner) {
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

// CheckSectionHeaderNames check section header names.
func CheckSectionHeaderNames(textPackID string) (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderNames(context, execution, textPackID)
	}
}

// CheckSectionHeaderDensity check section header density.
func CheckSectionHeaderDensity(textPackID string) (scanner runtimebinding.RepositoryScanner) {
	return func(
		context runner.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderDensity(context, execution, textPackID)
	}
}

// CheckSectionHeaders check section headers.
func CheckSectionHeaders(textPackID string) (scanner runtimebinding.RepositoryScanner) {
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
