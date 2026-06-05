package scan

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rules/golang/architecture"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/commandrun"
)

func scanGoArchitecture(
	context runner.Context,
	_ contract.RepositoryScanExecution,
) (result contract.ExecutionResult, err error) {
	modulePath, err := runGoList(context, "-m", "-f", "{{.Path}}")
	if err != nil {
		return contract.ExecutionResult{Output: modulePath}, err
	}

	packageList, err := runGoList(context, "-json", "./...")
	if err != nil {
		return contract.ExecutionResult{Output: packageList}, err
	}

	goConfig, err := decodeGoPackConfig(context)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	return architecture.CheckImports(modulePath, packageList, goConfig.Architecture)
}

func runGoList(context runner.Context, arguments ...string) (output string, err error) {
	return commandrun.Output(
		context.RepoRoot,
		context.GoEnvironment,
		"go",
		append([]string{"list"}, arguments...)...,
	)
}
