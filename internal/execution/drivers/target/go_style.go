package target

import (
	"errors"
	"fmt"
	"path/filepath"

	"ciphera/tools/internal/checks/golang"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/style"
)

// CheckGoStyle check go style.
func CheckGoStyle(goPackID string, goLanguage string) (check driverkit.TargetCheck) {
	return func(context execution.Context, job style.Job) (style.ExecutionResult, error) {
		return runGoStyleCheck(context, job, goPackID, goLanguage)
	}
}

func runGoStyleCheck(
	context execution.Context,
	job style.Job,
	goPackID string,
	goLanguage string,
) (result style.ExecutionResult, err error) {
	execution, found := job.(style.TargetCheckJob)
	if !found {
		return style.ExecutionResult{}, fmt.Errorf("go style check received empty job")
	}

	targets, err := goTargets(context, execution.Targets, goLanguage)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(context, goPackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	diagnostics := make([]style.Diagnostic, 0)
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
		joined = errors.Join(joined, err)
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, joined
}
