package target

import (
	"errors"
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
	return func(context runner.Context, spec style.ExecutionSpec) (style.ExecutionResult, error) {
		return runGolangci(context, spec, goPackID, golangciLintToolID, goimportsToolID, goLanguage)
	}
}

func runGolangci(
	context runner.Context,
	spec style.ExecutionSpec,
	goPackID string,
	golangciLintToolID string,
	goimportsToolID string,
	goLanguage string,
) (result style.ExecutionResult, err error) {
	if _, found := spec.TargetCommandExecution(); !found {
		return style.ExecutionResult{}, errEmptyTargetAction("golangci")
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
		workDir := targetWorkDir(context.RepoRoot, target)
		output, err := runGoFormatChecks(
			context,
			workDir,
			target.FormatPaths,
			localPrefix,
			goimportsToolID,
		)
		if err != nil {
			commandrun.AppendOutput(&builder, output)
			joined = errors.Join(joined, err)
			continue
		}

		output, err = commandrun.ToolByID(
			context,
			workDir,
			golangciLintToolID,
			"run",
			"./...",
		)
		commandrun.AppendOutput(&builder, output)
		joined = errors.Join(joined, err)
	}

	return style.ExecutionResult{Output: strings.TrimSpace(builder.String())}, joined
}

/* ---------------------------------------- Format Checks --------------------------------------- */

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

	if output, err = commandrun.Output(
		workDir,
		context.GoEnvironment,
		"gofmt",
		append([]string{"-l"}, paths...)...,
	); err != nil {
		return output, err
	}

	if strings.TrimSpace(output) != "" {
		return "Go files require gofmt formatting:\n" + strings.TrimSpace(output),
			errors.New("gofmt formatting required")
	}

	if output, err = commandrun.ToolByID(
		context,
		workDir,
		goimportsToolID,
		append([]string{"-l", "-local", localPrefix}, paths...)...,
	); err != nil {
		return output, err
	}

	if strings.TrimSpace(output) != "" {
		return "Go files require goimports formatting:\n" + strings.TrimSpace(output),
			errors.New("goimports formatting required")
	}

	return "", nil
}
