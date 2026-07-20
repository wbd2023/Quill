package target

import (
	"context"
	"strings"

	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/commandrun"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
)

/* --------------------------------------- Lint Execution --------------------------------------- */

// RunGolangci run golangci.
func RunGolangci(
	goPackID string,
	golangciLintToolID string,
	goimportsToolID string,
	goLanguage string,
) (command driverkit.TargetCommand) {
	return func(ctx context.Context, context execution.RunContext,
		job style.Job) (style.ExecutionResult, error) {
		return runGolangci(ctx, context, job, goPackID, golangciLintToolID, goimportsToolID,
			goLanguage)
	}
}

func runGolangci(
	ctx context.Context,
	context execution.RunContext,
	job style.Job,
	goPackID string,
	golangciLintToolID string,
	goimportsToolID string,
	goLanguage string,
) (result style.ExecutionResult, err error) {
	execution, found := job.(style.TargetCommandJob)
	if !found {
		return style.ExecutionResult{}, errEmptyTargetAction("golangci")
	}

	targets, err := goTargets(context, execution.Targets, goLanguage)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(context, goPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	var diagnostics []style.Diagnostic
	localPrefix := joinGoLocalImportPrefixes(goConfig.LocalImportPrefixes)
	for _, target := range targets {
		workDir := targetWorkDir(context.RepoRoot, target)
		output, err := runGoFormatChecks(
			ctx,
			context,
			workDir,
			target.FormatPaths,
			localPrefix,
			goimportsToolID,
		)
		if err != nil {
			return style.ExecutionResult{}, err
		}
		diagnostics = appendDiagnostics(diagnostics, output, "go/format")
		output, err = runGolangciLint(
			ctx,
			context,
			workDir,
			golangciLintToolID,
		)
		if err != nil {
			return style.ExecutionResult{}, err
		}
		diagnostics = appendDiagnostics(diagnostics, output, "go/lint")
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, nil
}

/* ---------------------------------------- Format Checks --------------------------------------- */

// runGolangciLint runs golangci-lint and returns its output. golangci-lint exits non-zero when it
// finds issues; that output is findings (data), not an operational error. Only command-execution
// failures (tool missing, timeout) produce a non-nil error.
func runGolangciLint(
	ctx context.Context,
	context execution.RunContext,
	workDir string,
	golangciLintToolID string,
) (output string, err error) {
	result, err := commandrun.ToolByID(
		ctx,
		context,
		workDir,
		golangciLintToolID,
		"run",
		"./...",
	)
	if err == nil {
		return "", nil
	}

	if result.ExitCode == 1 {
		return result.Output, nil
	}

	return result.Output, err
}

func runGoFormatChecks(
	ctx context.Context,
	context execution.RunContext,
	workDir string,
	paths []string,
	localPrefix string,
	goimportsToolID string,
) (output string, err error) {
	if len(paths) == 0 {
		return "", nil
	}

	gofmtResult, err := commandrun.Output(
		ctx,
		workDir,
		context.GoEnvironment,
		"gofmt",
		append([]string{"-l"}, paths...)...,
	)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(gofmtResult.Output) != "" {
		output = "Go files require gofmt formatting:\n" + strings.TrimSpace(gofmtResult.Output)
	}

	goimportsResult, err := commandrun.ToolByID(
		ctx,
		context,
		workDir,
		goimportsToolID,
		append([]string{"-l", "-local", localPrefix}, paths...)...,
	)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(goimportsResult.Output) != "" {
		if output != "" {
			output += "\n"
		}
		output += "Go files require goimports formatting:\n" +
			strings.TrimSpace(goimportsResult.Output)
	}

	return output, nil
}
