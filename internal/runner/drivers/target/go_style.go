package target

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rules/golang"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/commandrun"
)

func runGoStyleCheck(
	context runner.Context,
	spec contract.ExecutionSpec,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.TargetCheckExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("go style check received empty spec")
	}

	targets, err := goTargets(context, spec)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(context)
	if err != nil {
		return contract.ExecutionResult{}, err
	}

	diagnostics := make([]contract.Diagnostic, 0)
	var builder strings.Builder
	var joined error
	for _, target := range targets {
		if len(target.CheckPaths) == 0 {
			joined = errors.Join(
				joined,
				fmt.Errorf("go style target %q has no check paths", target.Name),
			)
			continue
		}

		workDir := targetWorkDir(context.RepoRoot, target)
		roots := make([]string, 0, len(target.CheckPaths))
		for _, checkPath := range target.CheckPaths {
			roots = append(roots, filepath.Join(workDir, checkPath))
		}

		styleResult, err := golang.CheckDirectories(
			context.RepoRoot,
			roots,
			context.Profile.Repository,
			context.Profile.PathRoles,
			goConfig,
			execution.Check,
		)
		diagnostics = append(diagnostics, styleResult.Diagnostics...)
		commandrun.AppendOutput(&builder, styleResult.Output)
		joined = errors.Join(joined, err)
	}

	return contract.ExecutionResult{
		Diagnostics: diagnostics,
		Output:      strings.TrimSpace(builder.String()),
	}, joined
}
