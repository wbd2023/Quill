package executors

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rules/golang"
	"ciphera/tools/internal/runner"
)

func runGoArchitectureCheck(
	context runner.Context,
) (result contract.ExecutionResult, err error) {
	modulePath, err := runGoList(context, "-m", "-f", "{{.Path}}")
	if err != nil {
		return contract.ExecutionResult{Output: modulePath}, err
	}

	packageList, err := runGoList(context, "-json", "./...")
	if err != nil {
		return contract.ExecutionResult{Output: packageList}, err
	}

	return golang.CheckArchitecture(modulePath, packageList, context.Policy.Go.Architecture)
}

func runGoList(context runner.Context, arguments ...string) (output string, err error) {
	return runCommandOutput(
		context.RepoRoot,
		context.GoEnvironment,
		"go",
		append([]string{"list"}, arguments...)...,
	)
}
