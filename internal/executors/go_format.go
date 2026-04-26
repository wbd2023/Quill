package executors

import (
	"errors"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

func runGoFormat(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error) {
	if _, found := spec.BackendCommandExecution(); !found {
		return contract.ExecutionResult{}, errEmptyBackendAction("go format")
	}

	backends, err := goLanguageBackends(context, spec)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	var builder strings.Builder
	var joined error
	for _, backend := range backends {
		if len(backend.FormatPaths) == 0 {
			continue
		}

		workdir := languageBackendWorkdir(context.RepoRoot, backend)
		output, err := runCommandOutput(
			workdir,
			context.GoEnvironment,
			"gofmt",
			append([]string{"-w"}, backend.FormatPaths...)...,
		)
		if err != nil {
			appendExecutorOutput(&builder, output)
			joined = errors.Join(joined, err)
			continue
		}

		output, err = runToolByID(
			context,
			workdir,
			rulepack.ToolGoimports,
			append(
				[]string{"-w", "-local", context.Policy.Imports.LocalPrefix},
				backend.FormatPaths...,
			)...,
		)
		appendExecutorOutput(&builder, output)
		joined = errors.Join(joined, err)
	}

	return contract.ExecutionResult{Output: strings.TrimSpace(builder.String())}, joined
}
