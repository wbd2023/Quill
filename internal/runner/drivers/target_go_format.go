package drivers

import (
	"errors"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/runner"
)

func runGoFormat(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error) {
	if _, found := spec.TargetCommandExecution(); !found {
		return contract.ExecutionResult{}, errEmptyTargetAction("go format")
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
		if len(target.FormatPaths) == 0 {
			continue
		}

		workDir := targetWorkDir(context.RepoRoot, target)
		output, err := runCommandOutput(
			workDir,
			context.GoEnvironment,
			"gofmt",
			append([]string{"-w"}, target.FormatPaths...)...,
		)
		if err != nil {
			appendDriverOutput(&builder, output)
			joined = errors.Join(joined, err)
			continue
		}

		output, err = runToolByID(
			context,
			workDir,
			builtin.ToolGoimports,
			append(
				[]string{"-w", "-local", localPrefix},
				target.FormatPaths...,
			)...,
		)
		appendDriverOutput(&builder, output)
		joined = errors.Join(joined, err)
	}

	return contract.ExecutionResult{Output: strings.TrimSpace(builder.String())}, joined
}
