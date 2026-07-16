package target

import (
	"errors"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/commandrun"
	"ciphera/tools/internal/execution/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

// RunGoFormat run go format.
func RunGoFormat(
	goPackID string,
	goimportsToolID string,
	goLanguage string,
) (command runtimebinding.TargetCommand) {
	return func(context execution.Context, job style.Job) (style.ExecutionResult, error) {
		return runGoFormat(context, job, goPackID, goimportsToolID, goLanguage)
	}
}

func runGoFormat(
	context execution.Context,
	job style.Job,
	goPackID string,
	goimportsToolID string,
	goLanguage string,
) (result style.ExecutionResult, err error) {
	execution, found := job.(style.TargetCommandJob)
	if !found {
		return style.ExecutionResult{}, errEmptyTargetAction("go format")
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
	var joined error
	localPrefix := joinGoLocalImportPrefixes(goConfig.LocalImportPrefixes)
	for _, target := range targets {
		if len(target.FormatPaths) == 0 {
			continue
		}

		workDir := targetWorkDir(context.RepoRoot, target)
		commandResult, err := commandrun.Output(
			workDir,
			context.GoEnvironment,
			"gofmt",
			append([]string{"-w"}, target.FormatPaths...)...,
		)
		if err != nil {
			diagnostics = appendDiagnostics(diagnostics, commandResult.Output, "go/format")
			joined = errors.Join(joined, err)
			continue
		}

		commandResult, err = commandrun.ToolByID(
			context,
			workDir,
			goimportsToolID,
			append(
				[]string{"-w", "-local", localPrefix},
				target.FormatPaths...,
			)...,
		)
		diagnostics = appendDiagnostics(diagnostics, commandResult.Output, "go/format")
		joined = errors.Join(joined, err)
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, joined
}
