package executors

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/report"
	"ciphera/tools/internal/rulepack"
	repostyle "ciphera/tools/internal/rules/repo"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runtime"
)

/* -------------------------------------------- Types ------------------------------------------- */

type qualityMakefileSurface struct {
	Variables map[string]string
	Targets   map[string]qualityMakefileTarget
}

type qualityMakefileTarget struct {
	Recipes []string
}

/* ------------------------------------ Control Plane Checks ------------------------------------ */

func controlPlaneExecutor(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]runtime.ToolStatus,
) (output string, err error) {
	switch spec.Check {
	case rulepack.ControlPlaneCheckEnforcementLevels:
		return checkEnforcementLevels()

	case rulepack.ControlPlaneCheckGlobalExclusions:
		return checkGlobalExclusions(context.Policy.Repository)

	case rulepack.ControlPlaneCheckQualityTargets:
		return checkQualityTargets(context.RepoRoot, context.Policy.ControlPlane)

	default:
		return "", fmt.Errorf("unknown control-plane check %q", spec.Check)
	}
}

func checkEnforcementLevels() (output string, err error) {
	requiredRule := contract.Rule{Level: contract.LevelRequired}
	recommendationRule := contract.Rule{Level: contract.LevelRecommendation}
	violation := errors.New("violation")

	switch runner.CheckStatus(requiredRule, violation, false) {
	case report.CheckStatusFail:
	default:
		return "required rules must fail on violations", errViolationsFound
	}

	switch runner.CheckStatus(recommendationRule, violation, false) {
	case report.CheckStatusWarn:
	default:
		return "recommendation rules must warn by default", errViolationsFound
	}

	switch runner.CheckStatus(recommendationRule, violation, true) {
	case report.CheckStatusFail:
	default:
		return "strict recommendations must fail on recommendation violations", errViolationsFound
	}

	return "", nil
}

func checkGlobalExclusions(repository profile.RepositoryConfig) (output string, err error) {
	if err = repostyle.ValidateCollectorPolicy(repository); err != nil {
		return err.Error(), errViolationsFound
	}

	return "", nil
}

func checkQualityTargets(
	repoRoot string,
	controlPlane profile.ControlPlaneConfig,
) (output string, err error) {
	contents, err := os.ReadFile(filepath.Join(repoRoot, controlPlane.QualityFile))
	if err != nil {
		return "", err
	}

	surface := parseQualityMakefileSurface(string(contents))
	for _, contract := range controlPlane.VariableContracts {
		actualValue, found := surface.Variables[contract.Name]
		if !found {
			return fmt.Sprintf(
				"%s is missing required variable: %s",
				controlPlane.QualityFile,
				contract.Name,
			), errViolationsFound
		}

		if actualValue == contract.Value {
			continue
		}

		return fmt.Sprintf(
			"%s variable %s must be %q, got %q",
			controlPlane.QualityFile,
			contract.Name,
			contract.Value,
			actualValue,
		), errViolationsFound
	}

	for _, contract := range controlPlane.TargetContracts {
		target, found := surface.Targets[contract.Name]
		if !found {
			return fmt.Sprintf(
				"%s is missing required target: %s",
				controlPlane.QualityFile,
				contract.Name,
			), errViolationsFound
		}

		if hasRecipeLine(target.Recipes, contract.RecipeLine) {
			continue
		}

		return fmt.Sprintf(
			"%s target %s is missing recipe line: %s",
			controlPlane.QualityFile,
			contract.Name,
			contract.RecipeLine,
		), errViolationsFound
	}

	return "", nil
}

/* -------------------------------------- Makefile Parsing -------------------------------------- */

func parseQualityMakefileSurface(contents string) (surface qualityMakefileSurface) {
	surface = qualityMakefileSurface{
		Variables: make(map[string]string),
		Targets:   make(map[string]qualityMakefileTarget),
	}

	currentTarget := ""
	for _, line := range strings.Split(contents, "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "", strings.HasPrefix(trimmed, "#"):
			currentTarget = ""

		case strings.HasPrefix(line, "\t"):
			if currentTarget == "" {
				continue
			}

			target := surface.Targets[currentTarget]
			target.Recipes = append(target.Recipes, strings.TrimSpace(line))
			surface.Targets[currentTarget] = target

		case strings.Contains(trimmed, "=") && !strings.Contains(trimmed, ":"):
			name, value, _ := strings.Cut(trimmed, "=")
			surface.Variables[strings.TrimSpace(name)] = strings.TrimSpace(value)
			currentTarget = ""

		case strings.Contains(trimmed, ":"):
			targetName, _, _ := strings.Cut(trimmed, ":")
			targetName = strings.TrimSpace(targetName)
			surface.Targets[targetName] = qualityMakefileTarget{}
			currentTarget = targetName

		default:
			currentTarget = ""
		}
	}

	return surface
}

/* --------------------------------------- Recipe Matching -------------------------------------- */

func hasRecipeLine(lines []string, expected string) (found bool) {
	for _, line := range lines {
		if strings.TrimSpace(line) == expected {
			return true
		}
	}

	return false
}
