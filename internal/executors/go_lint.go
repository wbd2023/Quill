package executors

import (
	"errors"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

func runGolangci(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error) {
	if _, found := spec.BackendCommandExecution(); !found {
		return contract.ExecutionResult{}, errEmptyBackendAction("golangci")
	}

	backends, err := goLanguageBackends(context, spec)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	var builder strings.Builder
	var joined error
	for _, backend := range backends {
		workDir := languageBackendWorkDir(context.RepoRoot, backend)
		output, err := runGoFormatChecks(context, workDir, backend.FormatPaths)
		if err != nil {
			appendExecutorOutput(&builder, output)
			joined = errors.Join(joined, err)
			continue
		}

		output, err = runToolByID(
			context,
			workDir,
			rulepack.ToolGolangciLint,
			"run",
			"./...",
		)
		appendExecutorOutput(&builder, output)
		joined = errors.Join(joined, err)
	}

	return contract.ExecutionResult{Output: strings.TrimSpace(builder.String())}, joined
}

func runGoFormatChecks(
	context runner.Context,
	workDir string,
	paths []string,
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

	localPrefix := joinGoLocalImportPrefixes(context.Policy.Go.LocalImportPrefixes)
	if output, err = runToolByID(
		context,
		workDir,
		rulepack.ToolGoimports,
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
