package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/policy"
	projectrules "ciphera/tools/internal/rules/project"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/toolchain"
)

var errViolationsFound = errors.New("violations found")

/* --------------------------------------- Project Checks --------------------------------------- */

func projectDriver(
	context runner.Context,
	spec contract.ExecutionSpec,
	_ map[string]toolchain.Status,
) (result contract.ExecutionResult, err error) {
	execution, found := spec.ProjectExecution()
	if !found {
		return contract.ExecutionResult{}, fmt.Errorf("project driver received empty spec")
	}

	var output string
	switch execution.Check {
	case builtin.ProjectCheckEnforcementLevels:
		output, err = checkEnforcementLevels()

	case builtin.ProjectCheckExcludedDirectories:
		output, err = checkExcludedDirectories(context.Profile.Repository)

	case builtin.ProjectCheckCommands:
		projectConfig, decodeErr := decodeProjectConfig(context)
		if decodeErr != nil {
			return contract.ExecutionResult{}, decodeErr
		}

		output, err = checkCommands(context.RepoRoot, projectConfig.Commands)

	default:
		return contract.ExecutionResult{}, fmt.Errorf(
			"unknown project check %q",
			execution.Check,
		)
	}

	return contract.ExecutionResult{Output: output}, err
}

func checkEnforcementLevels() (output string, err error) {
	requiredRule := contract.Rule{Enforcement: contract.EnforcementRequired}
	recommendationRule := contract.Rule{Enforcement: contract.EnforcementRecommendation}
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

func checkExcludedDirectories(repository policy.RepositoryConfig) (output string, err error) {
	if err = filewalk.ValidateCollectorPolicy(repository); err != nil {
		return err.Error(), errViolationsFound
	}

	return "", nil
}

func checkCommands(
	repositoryRoot string,
	commands projectrules.CommandsConfig,
) (output string, err error) {
	switch commands.Runner {
	case projectrules.CommandsRunnerMake:
		return checkMakeCommands(repositoryRoot, commands)
	default:
		return "", fmt.Errorf("unsupported quality commands runner %q", commands.Runner)
	}
}

func checkMakeCommands(
	repositoryRoot string,
	commands projectrules.CommandsConfig,
) (output string, err error) {
	contents, err := os.ReadFile(filepath.Join(repositoryRoot, commands.Make.Path))
	if err != nil {
		return "", err
	}

	makefile := parseMakefileSurface(string(contents))
	for _, variable := range commands.Make.RequiredVariables {
		actual, found := makefile.Variables[variable.Name]
		if !found {
			return fmt.Sprintf(
				"%s is missing required variable: %s",
				commands.Make.Path,
				variable.Name,
			), errViolationsFound
		}

		if actual == variable.Value {
			continue
		}

		return fmt.Sprintf(
			"%s variable %s must be %q, got %q",
			commands.Make.Path,
			variable.Name,
			variable.Value,
			actual,
		), errViolationsFound
	}

	for _, requiredTarget := range commands.Make.RequiredTargets {
		target, found := makefile.Targets[requiredTarget.Name]
		if !found {
			return fmt.Sprintf(
				"%s is missing required target: %s",
				commands.Make.Path,
				requiredTarget.Name,
			), errViolationsFound
		}

		if hasRecipeLine(target.Recipes, requiredTarget.RecipeLine) {
			continue
		}

		return fmt.Sprintf(
			"%s target %s is missing recipe line: %s",
			commands.Make.Path,
			requiredTarget.Name,
			requiredTarget.RecipeLine,
		), errViolationsFound
	}

	return "", nil
}
