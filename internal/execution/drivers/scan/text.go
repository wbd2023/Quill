package scan

import (
	"ciphera/tools/internal/checks/text"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
)

// CheckASCII scans for non-ASCII characters in text files.
func CheckASCII() (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.Context,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return text.CheckASCII(context.RepoRoot, context.Profile.Repository, context.Scope)
	}
}

// CheckExceptionMarkers check exception markers.
func CheckExceptionMarkers() (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.Context,
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
func CheckLineLengths() (scanner driverkit.RepositoryScanner) {
	return scanLineLengths
}

// CheckMaintenanceMarkers check maintenance markers.
func CheckMaintenanceMarkers() (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.Context,
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
func CheckSectionHeaderNames(textPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderNames(context, execution, textPackID)
	}
}

// CheckSectionHeaderDensity check section header density.
func CheckSectionHeaderDensity(textPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderDensity(context, execution, textPackID)
	}
}

// CheckSectionHeaders check section headers.
func CheckSectionHeaders(textPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		context execution.Context,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaders(context, execution, textPackID)
	}
}

func scanLineLengths(
	context execution.Context,
	scanExec style.RepositoryScanExecution,
) (result style.ExecutionResult, err error) {
	files, err := execution.CollectFileSetFiles(context, scanExec.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return text.CheckLineLengths(context.RepoRoot, files)
}
