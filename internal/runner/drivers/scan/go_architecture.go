package scan

import (
	"fmt"

	"ciphera/tools/internal/checks/golang/architecture"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/commandrun"
	"ciphera/tools/internal/style"
)

func scanGoArchitecture(
	context runner.Context,
	_ style.RepositoryScanExecution,
	goPackID string,
) (result style.ExecutionResult, err error) {
	modulePath, err := runGoList(context, "-m", "-f", "{{.Path}}")
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("go list module path: %w", err)
	}

	packageList, err := runGoList(context, "-json", "./...")
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("go list packages: %w", err)
	}

	goConfig, err := decodeGoPackConfig(context, goPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return architecture.CheckImports(modulePath, packageList, goConfig.Architecture)
}

func runGoList(context runner.Context, arguments ...string) (output string, err error) {
	result, err := commandrun.Output(
		context.RepoRoot,
		context.GoEnvironment,
		"go",
		append([]string{"list"}, arguments...)...,
	)
	return result.Output, err
}
