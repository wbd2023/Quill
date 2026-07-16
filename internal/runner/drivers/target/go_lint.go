package target

import (
	"strings"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/commandrun"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

/* --------------------------------------- Lint Execution --------------------------------------- */

// RunGolangci run golangci.
func RunGolangci(
	goPackID string,
	golangciLintToolID string,
	goimportsToolID string,
	goLanguage string,
) (command runtimebinding.TargetCommand) {
	return func(context runner.Context, job style.Job) (style.ExecutionResult, error) {
		return runGolangci(context, job, goPackID, golangciLintToolID, goimportsToolID, goLanguage)
	}
}

func runGolangci(
	context runner.Context,
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
	context runner.Context,
	workDir string,
	golangciLintToolID string,
) (output string, err error) {
	result, err := commandrun.ToolByID(
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
	context runner.Context,
	workDir string,
	paths []string,
	localPrefix string,
	goimportsToolID string,
) (output string, err error) {
	if len(paths) == 0 {
		return "", nil
	}

	gofmtResult, err := commandrun.Output(
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
