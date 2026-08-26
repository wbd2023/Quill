package bindings

import (
	"context"
	"errors"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/style"
)

// formatTargets is the Go go_format target command: it applies gofmt -w and goimports -w to each
// target's format paths. Formatting failures surface their output as diagnostics and do not abort
// the remaining targets.
func formatTargets(
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
	var joined error
	localPrefix := joinGoLocalImportPrefixes(goConfig.LocalImportPrefixes)
	for _, target := range targets {
		if len(target.FormatPaths) == 0 {
			continue
		}

		workDir := targetWorkDir(run.Root, target)
		commandResult, err := runGo(
			ctx,
			run,
			workDir,
			"gofmt",
			append([]string{"-w"}, target.FormatPaths...)...,
		)
		if err != nil {
			diagnostics = appendDiagnostics(diagnostics, commandResult.Output, "go/format")
			joined = errors.Join(joined, err)
			continue
		}

		commandResult, err = runTool(
			ctx,
			run,
			workDir,
			tool.Goimports,
			append([]string{"-w", "-local", localPrefix}, target.FormatPaths...)...,
		)
		diagnostics = appendDiagnostics(diagnostics, commandResult.Output, "go/format")
		joined = errors.Join(joined, err)
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, joined
}
