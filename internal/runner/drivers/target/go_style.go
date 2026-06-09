package target

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/checks/golang"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/runner/drivers/internal/commandrun"
	"ciphera/tools/internal/style"
)

func CheckGoStyle(goPackID string, goLanguage string) (check binding.TargetCheck) {
	return func(context runner.Context, spec style.ExecutionSpec) (style.ExecutionResult, error) {
		return runGoStyleCheck(context, spec, goPackID, goLanguage)
	}
}

func runGoStyleCheck(
	context runner.Context,
	spec style.ExecutionSpec,
	goPackID string,
	goLanguage string,
) (result style.ExecutionResult, err error) {
	execution, found := spec.TargetCheckExecution()
	if !found {
		return style.ExecutionResult{}, fmt.Errorf("go style check received empty spec")
	}

	targets, err := goTargets(context, spec, goLanguage)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(context, goPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	diagnostics := make([]style.Diagnostic, 0)
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

	return style.ExecutionResult{
		Diagnostics: diagnostics,
		Output:      strings.TrimSpace(builder.String()),
	}, joined
}
