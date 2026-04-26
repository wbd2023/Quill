package executors

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------ Control Plane Checks ------------------------------------ */

func controlPlaneExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.ControlPlaneExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("control-plane executor received empty spec")
	}

	var output string
	switch execution.Check {
	case rulepack.ControlPlaneCheckEnforcementLevels:
		output, err = checkEnforcementLevels()

	case rulepack.ControlPlaneCheckGlobalExclusions:
		output, err = checkGlobalExclusions(context.Policy.Repository)

	case rulepack.ControlPlaneCheckQualityTargets:
		output, err = checkQualityTargets(context.RepoRoot, context.Policy.ControlPlane)

	default:
		return contract.ExecutionResult{}, fmt.Errorf(
			"unknown control-plane check %q",
			execution.Check,
		)
	}

	return contract.ExecutionResult{Output: output}, err
}

func checkEnforcementLevels() (output string, err error) {
	requiredRule := contract.Rule{Level: contract.LevelRequired}
	recommendationRule := contract.Rule{Level: contract.LevelRecommendation}
	violation := errors.New("violation")

	switch runner.CheckStatus(requiredRule, violation, false) {
	case contract.CheckStatusFail:
	default:
		return "required rules must fail on violations", errViolationsFound
	}

	switch runner.CheckStatus(recommendationRule, violation, false) {
	case contract.CheckStatusWarn:
	default:
		return "recommendation rules must warn by default", errViolationsFound
	}

	switch runner.CheckStatus(recommendationRule, violation, true) {
	case contract.CheckStatusFail:
	default:
		return "strict recommendations must fail on recommendation violations", errViolationsFound
	}

	return "", nil
}

func checkGlobalExclusions(repository policy.RepositoryConfig) (output string, err error) {
	if err = filewalk.ValidateCollectorPolicy(repository); err != nil {
		return err.Error(), errViolationsFound
	}

	return "", nil
}

func checkQualityTargets(
	repoRoot string,
	controlPlane policy.ControlPlaneConfig,
) (output string, err error) {
	contents, err := os.ReadFile(filepath.Join(repoRoot, controlPlane.QualityFile))
	if err != nil {
		return "", err
	}

	surface := parseQualityMakefileSurface(string(contents))
	for _, variable := range controlPlane.VariableContracts {
		actual, found := surface.Variables[variable.Name]
		if !found {
			return fmt.Sprintf(
				"%s is missing required variable: %s",
				controlPlane.QualityFile,
				variable.Name,
			), errViolationsFound
		}

		if actual == variable.Value {
			continue
		}

		return fmt.Sprintf(
			"%s variable %s must be %q, got %q",
			controlPlane.QualityFile,
			variable.Name,
			variable.Value,
			actual,
		), errViolationsFound
	}

	for _, requiredTarget := range controlPlane.TargetContracts {
		target, found := surface.Targets[requiredTarget.Name]
		if !found {
			return fmt.Sprintf(
				"%s is missing required target: %s",
				controlPlane.QualityFile,
				requiredTarget.Name,
			), errViolationsFound
		}

		if hasRecipeLine(target.Recipes, requiredTarget.RecipeLine) {
			continue
		}

		return fmt.Sprintf(
			"%s target %s is missing recipe line: %s",
			controlPlane.QualityFile,
			requiredTarget.Name,
			requiredTarget.RecipeLine,
		), errViolationsFound
	}

	return "", nil
}
