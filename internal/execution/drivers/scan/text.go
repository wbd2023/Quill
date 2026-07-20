package scan

import (
	"context"

	"github.com/wbd2023/Quill/internal/checks/text"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
)

/* ---------------------------------------- Text Scanners --------------------------------------- */

// CheckASCII scans for non-ASCII characters in text files.
func CheckASCII() (scanner driverkit.RepositoryScanner) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return text.CheckASCII(context.RepoRoot, context.Profile.Repository, context.Scope)
	}
}

// CheckExceptionMarkers check exception markers.
func CheckExceptionMarkers() (scanner driverkit.RepositoryScanner) {
	return func(
		_ context.Context,
		context execution.RunContext,
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
		_ context.Context,
		context execution.RunContext,
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
		ctx context.Context,
		context execution.RunContext,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderNames(ctx, context, execution, textPackID)
	}
}

// CheckSectionHeaderDensity check section header density.
func CheckSectionHeaderDensity(textPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaderDensity(ctx, context, execution, textPackID)
	}
}

// CheckSectionHeaders check section headers.
func CheckSectionHeaders(textPackID string) (scanner driverkit.RepositoryScanner) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		execution style.RepositoryScanExecution,
	) (style.ExecutionResult, error) {
		return scanSectionHeaders(ctx, context, execution, textPackID)
	}
}

func scanLineLengths(
	_ context.Context,
	context execution.RunContext,
	scanExec style.RepositoryScanExecution,
) (result style.ExecutionResult, err error) {
	files, err := execution.CollectFileSetFiles(context, scanExec.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return text.CheckLineLengths(context.RepoRoot, files)
}
