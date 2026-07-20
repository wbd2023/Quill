package profile

import (
	"context"
	"fmt"

	"github.com/wbd2023/Quill/internal/checks/project"
	"github.com/wbd2023/Quill/internal/execution"
	"github.com/wbd2023/Quill/internal/execution/drivers/internal/driverkit"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

/* --------------------------------------- Profile Driver --------------------------------------- */

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

// CheckDriver returns the profile driver for check execution.
func CheckDriver(checks driverkit.ProfileChecks) (driver execution.Driver) {
	return func(
		ctx context.Context,
		context execution.RunContext,
		job style.Job,
		_ toolchain.StatusMap,
	) (result style.ExecutionResult, err error) {
		execution, found := job.(style.ProfileExecution)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf("profile driver received empty job")
		}

		check, found := checks.Lookup(execution.Check)
		if !found {
			return style.ExecutionResult{}, fmt.Errorf(
				"unknown profile check %q",
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
		message, err := project.CheckEnforcementLevels()
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
		message, err := project.CheckExcludedDirectories(context.Profile.Repository)
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

		message, err := project.CheckCommands(context.RepoRoot, projectConfig.Commands)
		return enforcementResult(message), err
	}
}
