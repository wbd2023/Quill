package drivers

import (
	"errors"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/runner"
)

/* --------------------------------------- Lint Execution --------------------------------------- */

func runGolangci(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error) {
	if _, found := spec.TargetCommandExecution(); !found {
		return contract.ExecutionResult{}, errEmptyTargetAction("golangci")
	}

	targets, err := goTargets(context, spec)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(context)
	if err != nil {
		return contract.ExecutionResult{}, err
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
		)
		if err != nil {
			appendExecutorOutput(&builder, output)
			joined = errors.Join(joined, err)
			continue
		}

		output, err = runToolByID(
			context,
			workDir,
			builtin.ToolGolangciLint,
			"run",
			"./...",
		)
		appendExecutorOutput(&builder, output)
		joined = errors.Join(joined, err)
	}

	return contract.ExecutionResult{Output: strings.TrimSpace(builder.String())}, joined
}

/* ---------------------------------------- Format Checks --------------------------------------- */

func runGoFormatChecks(
	context runner.Context,
	workDir string,
	paths []string,
	localPrefix string,
) (output string, err error) {
	if len(paths) == 0 {
		return "", nil
	}

	if output, err = runCommandOutput(
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

	if output, err = runToolByID(
		context,
		workDir,
		builtin.ToolGoimports,
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
