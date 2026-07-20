package project

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wbd2023/Quill/internal/checks/projectpolicy"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/filewalk"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
)

/* --------------------------------------- Project Checks --------------------------------------- */

var errCheckStatusMisconfigured = errors.New("check status classification is misconfigured")

// CheckEnforcementLevels checks the required and recommendation status classifications.
func CheckEnforcementLevels() (message string, err error) {
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

// CheckExcludedDirectories checks whether repository exclusions form a valid collector policy.
func CheckExcludedDirectories(
	repository policy.RepositoryConfig,
) (message string, err error) {
	if err = filewalk.ValidateCollectorPolicy(filewalk.WalkConfig{
		ExcludedDirectories: repository.ExcludedDirectories,
		GeneratedMarker:     repository.GeneratedMarker,
	}); err != nil {
		return err.Error(), err
	}

	return "", nil
}

// CheckCommands checks the configured repository quality-command surface.
func CheckCommands(
	repositoryRoot string,
	commands projectpolicy.CommandsConfig,
) (message string, err error) {
	switch commands.Runner {
	case projectpolicy.CommandsRunnerMake:
		return checkMakeCommands(repositoryRoot, commands)
	default:
		return "", fmt.Errorf("unsupported quality commands runner %q", commands.Runner)
	}
}

/* --------------------------------------- Makefile Checks -------------------------------------- */

func checkMakeCommands(
	repositoryRoot string,
	commands projectpolicy.CommandsConfig,
) (message string, err error) {
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
