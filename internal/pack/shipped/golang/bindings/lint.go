package bindings

import (
	"context"
	"strings"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/style"
)

/* ----------------------------------- Golangci Target Command ---------------------------------- */

// lintTargets is the Go golangci target command: for each target it runs the format checks
// (gofmt/goimports) and then golangci-lint, surfacing findings as diagnostics.
func lintTargets(
	ctx context.Context,
	run execution.RunContext,
	job style.TargetCommandJob,
) (result style.ExecutionResult, err error) {
	targets, err := goTargets(run, job.Targets, golang.Language)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(run, golang.PackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	var diagnostics []style.Diagnostic
	localPrefix := joinGoLocalImportPrefixes(goConfig.LocalImportPrefixes)
	for _, target := range targets {
		workDir := targetWorkDir(run.Root, target)
		output, err := runGoFormatChecks(
			ctx,
			run,
			workDir,
			target.FormatPaths,
			localPrefix,
			tool.Goimports,
		)
		if err != nil {
			return style.ExecutionResult{}, err
		}
		diagnostics = appendDiagnostics(diagnostics, output, "go/format")

		output, err = runGolangciLint(ctx, run, workDir, tool.GolangciLint)
		if err != nil {
			return style.ExecutionResult{}, err
		}
		diagnostics = appendDiagnostics(diagnostics, output, "go/lint")
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, nil
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// runGolangciLint runs golangci-lint and returns its output. golangci-lint exits non-zero when it
// finds issues; that output is findings (data), not an operational error. Only command-execution
// failures (tool missing, timeout) produce a non-nil error.
func runGolangciLint(
	ctx context.Context,
	run execution.RunContext,
	workDir string,
	toolID string,
) (output string, err error) {
	result, err := runTool(ctx, run, workDir, toolID, "run", "./...")
	if err == nil {
		return "", nil
	}

	if result.ExitCode == 1 {
		return result.Output, nil
	}

	return result.Output, err
}

// runGoFormatChecks runs gofmt -l and goimports -l over the given paths and returns combined
// findings output. An empty result means no formatting issues were found.
func runGoFormatChecks(
	ctx context.Context,
	run execution.RunContext,
	workDir string,
	paths []string,
	localPrefix string,
	toolID string,
) (output string, err error) {
	if len(paths) == 0 {
		return "", nil
	}

	gofmtResult, err := runGo(ctx, run, workDir, "gofmt", append([]string{"-l"}, paths...)...)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(gofmtResult.Output) != "" {
		output = "Go files require gofmt formatting:\n" + strings.TrimSpace(gofmtResult.Output)
	}

	goimportsResult, err := runTool(
		ctx,
		run,
		workDir,
		toolID,
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
