package bindings

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	checks "github.com/wbd2023/quill/internal/checks/golang"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/golang"
	"github.com/wbd2023/quill/internal/style"
)

// styleCheck returns the Go target check closure for a single Check selector. Each TargetCheck*
// identity is bound to exactly one selector; the closure runs that selector over every target's
// check paths.
func styleCheck(selector checks.Check) (check drivers.TargetCheck) {
	return func(
		_ context.Context,
		run execution.RunContext,
		job style.TargetCheckJob,
	) (result style.ExecutionResult, err error) {
		return runStyleCheck(run, job, selector)
	}
}

func runStyleCheck(
	run execution.RunContext,
	job style.TargetCheckJob,
	selector checks.Check,
) (result style.ExecutionResult, err error) {
	targets, err := goTargets(run, job.Targets, golang.Language)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	goConfig, err := decodeGoConfig(run, golang.PackID)
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

		workDir := targetWorkDir(run.RepoRoot, target)
		roots := make([]string, 0, len(target.CheckPaths))
		for _, checkPath := range target.CheckPaths {
			roots = append(roots, filepath.Join(workDir, checkPath))
		}

		styleResult, err := checks.CheckDirectories(
			run.RepoRoot,
			roots,
			run.Profile.Repository,
			run.Profile.PathRoles,
			goConfig,
			selector,
		)
		diagnostics = append(diagnostics, styleResult.Diagnostics...)
		joined = errors.Join(joined, err)
	}

	return style.ExecutionResult{Diagnostics: diagnostics}, joined
}
