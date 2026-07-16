package target

import (
	"fmt"
	"path/filepath"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/policy"
)

func goTargets(
	context execution.Context,
	names []string,
	goLanguage string,
) (targets []policy.TargetConfig, err error) {
	for _, name := range names {
		target, err := goTarget(context.Profile, name, goLanguage)
		if err != nil {
			return nil, err
		}

		if !context.Profile.Repository.HasScopeOverlap(context.Scope, target.Scope) {
			continue
		}

		targets = append(targets, target)
	}

	return targets, nil
}

func goTarget(
	config policy.Config,
	name string,
	goLanguage string,
) (target policy.TargetConfig, err error) {
	target, found := config.Targets.Lookup(name)
	if !found {
		return policy.TargetConfig{}, fmt.Errorf("unknown Go target %q", name)
	}

	if target.Language != goLanguage {
		return policy.TargetConfig{}, fmt.Errorf(
			"target %q is %q, not go",
			name,
			target.Language,
		)
	}

	return target, nil
}

func targetWorkDir(
	repoRoot string,
	target policy.TargetConfig,
) (workDir string) {
	if target.WorkingDirectory == "" || target.WorkingDirectory == "." {
		return repoRoot
	}

	return filepath.Join(repoRoot, target.WorkingDirectory)
}

func errEmptyTargetAction(action string) (err error) {
	return fmt.Errorf("%s action received empty spec", action)
}
