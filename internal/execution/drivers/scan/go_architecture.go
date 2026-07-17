package scan

import (
	"context"
	"fmt"

	"ciphera/tools/internal/checks/golang/architecture"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/commandrun"
	"ciphera/tools/internal/style"
)

func scanGoArchitecture(
	ctx context.Context,
	context execution.RunContext,
	_ style.RepositoryScanExecution,
	goPackID string,
) (result style.ExecutionResult, err error) {
	modulePath, err := runGoList(ctx, context, "-m", "-f", "{{.Path}}")
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("go list module path: %w", err)
	}

	packageList, err := runGoList(ctx, context, "-json", "./...")
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("go list packages: %w", err)
	}

	goConfig, err := decodeGoPackConfig(context, goPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return architecture.CheckImports(modulePath, packageList, goConfig.Architecture)
}

func runGoList(ctx context.Context, context execution.RunContext,
	arguments ...string) (output string, err error) {
	result, err := commandrun.Output(
		ctx,
		context.RepoRoot,
		context.GoEnvironment,
		"go",
		append([]string{"list"}, arguments...)...,
	)
	return result.Output, err
}
