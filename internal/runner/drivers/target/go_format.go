package target

import (
	"errors"
	"strings"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/commandrun"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	"ciphera/tools/internal/style"
)

func RunGoFormat(
	goPackID string,
	goimportsToolID string,
	goLanguage string,
) (command runtimebinding.TargetCommand) {
	return func(context runner.Context, spec style.ExecutionSpec) (style.ExecutionResult, error) {
		return runGoFormat(context, spec, goPackID, goimportsToolID, goLanguage)
	}
}

func runGoFormat(
	context runner.Context,
	spec style.ExecutionSpec,
	goPackID string,
	goimportsToolID string,
	goLanguage string,
) (result style.ExecutionResult, err error) {
	if _, found := spec.TargetCommandExecution(); !found {
		return style.ExecutionResult{}, errEmptyTargetAction("go format")
	}

	targets, err := goTargets(context, spec, goLanguage)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(context, goPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	var builder strings.Builder
	var joined error
	localPrefix := joinGoLocalImportPrefixes(goConfig.LocalImportPrefixes)
	for _, target := range targets {
		if len(target.FormatPaths) == 0 {
			continue
		}

		workDir := targetWorkDir(context.RepoRoot, target)
		output, err := commandrun.Output(
			workDir,
			context.GoEnvironment,
			"gofmt",
			append([]string{"-w"}, target.FormatPaths...)...,
		)
		if err != nil {
			commandrun.AppendOutput(&builder, output)
			joined = errors.Join(joined, err)
			continue
		}

		output, err = commandrun.ToolByID(
			context,
			workDir,
			goimportsToolID,
			append(
				[]string{"-w", "-local", localPrefix},
				target.FormatPaths...,
			)...,
		)
		commandrun.AppendOutput(&builder, output)
		joined = errors.Join(joined, err)
	}

	return style.ExecutionResult{Output: strings.TrimSpace(builder.String())}, joined
}
