package bindings

import (
	"context"
	"fmt"

	"github.com/wbd2023/quill/internal/checks/golang/architecture"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	"github.com/wbd2023/quill/internal/style"
)

// scanArchitecture checks Go import boundaries against the configured architecture policy.
func scanArchitecture(
	ctx context.Context,
	run execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	modulePath, err := runGoList(ctx, run, "-m", "-f", "{{.Path}}")
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("go list module path: %w", err)
	}

	packagesJSON, err := runGoList(ctx, run, "-json", "./...")
	if err != nil {
		return style.ExecutionResult{}, fmt.Errorf("go list packages: %w", err)
	}

	goConfig, err := decodeGoConfig(run, golang.PackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return architecture.CheckImports(modulePath, packagesJSON, goConfig.Architecture)
}

func runGoList(
	ctx context.Context,
	run execution.RunContext,
	arguments ...string,
) (output string, err error) {
	result, err := runGo(
		ctx,
		run,
		run.Root,
		"go",
		append([]string{"list"}, arguments...)...,
	)
	return result.Output, err
}
