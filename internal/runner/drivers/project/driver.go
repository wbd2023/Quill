package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	projectpolicy "ciphera/tools/internal/checks/project/policy"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/runner/drivers/internal/binding"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

var errViolationsFound = errors.New("violations found")

/* --------------------------------------- Project Checks --------------------------------------- */

func projectDriver(checks binding.ProjectChecks) (driver runner.Driver) {
	return func(
		context runner.Context,
		spec style.ExecutionSpec,
		_ map[string]toolchain.Status,
	) (result style.ExecutionResult, err error) {
		execution, found := spec.ProjectExecution()
		if !found {
			return style.ExecutionResult{}, fmt.Errorf("project driver received empty spec")
		}

		check, found := checks.Lookup(execution.Check)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown project check %q",
				execution.Check,
			)
		}

		return check(context, execution)
	}
}

func CheckEnforcementLevels() (check binding.ProjectCheck) {
	return func(
		_ runner.Context,
		_ style.ProjectExecution,
	) (result style.ExecutionResult, err error) {
		output, err := checkEnforcementLevels()
		return style.ExecutionResult{Output: output}, err
	}
}

func CheckExcludedDirectories() (check binding.ProjectCheck) {
	return func(
		context runner.Context,
		_ style.ProjectExecution,
	) (result style.ExecutionResult, err error) {
		output, err := checkExcludedDirectories(context.Profile.Repository)
		return style.ExecutionResult{Output: output}, err
	}
}

func CheckCommands(projectPackID string) (check binding.ProjectCheck) {
	return func(
		context runner.Context,
		_ style.ProjectExecution,
	) (result style.ExecutionResult, err error) {
		projectConfig, err := decodeProjectConfig(context, projectPackID)
		if err != nil {
			return style.ExecutionResult{}, err
		}

		output, err := checkCommands(context.RepoRoot, projectConfig.Commands)
		return style.ExecutionResult{Output: output}, err
	}
}

func checkEnforcementLevels() (output string, err error) {
	requiredRule := style.Rule{Enforcement: style.EnforcementRequired}
	recommendationRule := style.Rule{Enforcement: style.EnforcementRecommendation}
	violation := errors.New("violation")

	switch runner.CheckStatus(requiredRule, violation, false) {
	case style.CheckStatusFail:
	default:
		return "required rules must fail on violations", errViolationsFound
	}

	switch runner.CheckStatus(recommendationRule, violation, false) {
	case style.CheckStatusWarn:
	default:
		return "recommendation rules must warn by default", errViolationsFound
	}

	switch runner.CheckStatus(recommendationRule, violation, true) {
	case style.CheckStatusFail:
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
	commands projectpolicy.CommandsConfig,
) (output string, err error) {
	switch commands.Runner {
	case projectpolicy.CommandsRunnerMake:
		return checkMakeCommands(repositoryRoot, commands)
	default:
		return "", fmt.Errorf("unsupported quality commands runner %q", commands.Runner)
	}
}

func checkMakeCommands(
	repositoryRoot string,
	commands projectpolicy.CommandsConfig,
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
