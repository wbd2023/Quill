package profile

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"ciphera/tools/internal/checks/projectpolicy"
	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/execution/drivers/internal/driverkit"
	"ciphera/tools/internal/filewalk"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

var errCheckStatusMisconfigured = errors.New("check status classification is misconfigured")

// enforcementResult converts a profile-check message into an ExecutionResult. A non-empty message
// is a finding; an empty message means the check passed.
func enforcementResult(message string) (result style.ExecutionResult) {
	if message == "" {
		return style.ExecutionResult{}
	}

	return style.ExecutionResult{
		Diagnostics: []style.Diagnostic{{
			Code:    "profile/enforcement",
			Message: message,
		}},
	}
}

/* --------------------------------------- Project Checks --------------------------------------- */

func profileDriver(checks driverkit.ProfileChecks) (driver execution.Executor) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		execution, found := job.(style.ProfileExecution)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf("project driver received empty job")
		}

		check, found := checks.Lookup(execution.Check)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown project check %q",
				execution.Check,
			)
		}

		return check(ctx, context, execution)
	}
}

// CheckEnforcementLevels check enforcement levels.
func CheckEnforcementLevels() (check driverkit.ProfileCheck) {
	return func(
		_ context.Context,
		_ execution.RunContext,
		_ style.ProfileExecution,
	) (result style.ExecutionResult, err error) {
		message, err := checkEnforcementLevels()
		return enforcementResult(message), err
	}
}

// CheckExcludedDirectories check excluded directories.
func CheckExcludedDirectories() (check driverkit.ProfileCheck) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.ProfileExecution,
	) (result style.ExecutionResult, err error) {
		message, err := checkExcludedDirectories(context.Profile.Repository)
		return enforcementResult(message), err
	}
}

// CheckCommands check commands.
func CheckCommands(profilePackID string) (check driverkit.ProfileCheck) {
	return func(
		_ context.Context,
		context execution.RunContext,
		_ style.ProfileExecution,
	) (result style.ExecutionResult, err error) {
		projectConfig, err := decodeProjectConfig(context, profilePackID)
		if err != nil {
			return style.ExecutionResult{}, err
		}

		message, err := checkCommands(context.RepoRoot, projectConfig.Commands)
		return enforcementResult(message), err
	}
}

func checkEnforcementLevels() (output string, err error) {
	requiredRule := style.Rule{Enforcement: style.EnforcementRequired}
	recommendationRule := style.Rule{Enforcement: style.EnforcementRecommendation}
	violations := style.ExecutionResult{
		Diagnostics: []style.Diagnostic{{Code: "self-test", Message: "violation"}},
	}

	switch execution.CheckStatus(requiredRule, violations, nil, false) {
	case style.CheckStatusFail:
	default:
		return "required rules must fail on violations", errCheckStatusMisconfigured
	}

	switch execution.CheckStatus(recommendationRule, violations, nil, false) {
	case style.CheckStatusWarn:
	default:
		return "recommendation rules must warn by default", errCheckStatusMisconfigured
	}

	switch execution.CheckStatus(recommendationRule, violations, nil, true) {
	case style.CheckStatusFail:
	default:
		return "strict recommendations must fail on recommendation violations",
			errCheckStatusMisconfigured
	}

	return "", nil
}

func checkExcludedDirectories(repository policy.RepositoryConfig) (output string, err error) {
	if err = filewalk.ValidateCollectorPolicy(filewalk.WalkConfig{
		ExcludedDirectories: repository.ExcludedDirectories,
		GeneratedMarker:     repository.GeneratedMarker,
	}); err != nil {
		return err.Error(), err
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
			), nil
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
		), nil
	}

	for _, requiredTarget := range commands.Make.RequiredTargets {
		target, found := makefile.Targets[requiredTarget.Name]
		if !found {
			return fmt.Sprintf(
				"%s is missing required target: %s",
				commands.Make.Path,
				requiredTarget.Name,
			), nil
		}

		if hasRecipeLine(target.Recipes, requiredTarget.RecipeLine) {
			continue
		}

		return fmt.Sprintf(
			"%s target %s is missing recipe line: %s",
			commands.Make.Path,
			requiredTarget.Name,
			requiredTarget.RecipeLine,
		), nil
	}

	return "", nil
}
