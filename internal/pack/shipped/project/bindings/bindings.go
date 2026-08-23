// Package bindings owns the Project Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Project execution identities (profile checks) to concrete
// check behaviour. The parent project package stays independent of the driver facade and check
// implementations.
package bindings

import (
	"context"
	"fmt"

	checks "github.com/wbd2023/quill/internal/checks/project"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/project"
	projectpolicy "github.com/wbd2023/quill/internal/pack/shipped/project/policy"
	"github.com/wbd2023/quill/internal/style"
)

// Register wires every Project execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	bindings.AddProfileCheck(project.PackID, project.CheckEnforcementLevels, checkEnforcementLevels)
	bindings.AddProfileCheck(
		project.PackID,
		project.CheckExcludedDirectories,
		checkExcludedDirectories,
	)
	bindings.AddProfileCheck(project.PackID, project.CheckCommands, checkCommands)
}

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

func checkEnforcementLevels(
	_ context.Context,
	_ execution.RunContext,
	_ style.ProfileCheck,
) (result style.ExecutionResult, err error) {
	message, err := checks.CheckEnforcementLevels()
	return enforcementResult(message), err
}

func checkExcludedDirectories(
	_ context.Context,
	context execution.RunContext,
	_ style.ProfileCheck,
) (result style.ExecutionResult, err error) {
	message, err := checks.CheckExcludedDirectories(context.Profile.Repository)
	return enforcementResult(message), err
}

func checkCommands(
	_ context.Context,
	context execution.RunContext,
	_ style.ProfileCheck,
) (result style.ExecutionResult, err error) {
	projectConfig, err := decodeProjectConfig(context, project.PackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	message, err := checks.CheckCommands(context.RepoRoot, projectConfig.Commands)
	return enforcementResult(message), err
}

func decodeProjectConfig(
	context execution.RunContext,
	packID string,
) (config projectpolicy.Config, err error) {
	pack, found := context.Profile.PackPolicies.Lookup(packID)
	if !found {
		return projectpolicy.Config{}, errMissingPackPolicy(packID)
	}

	return projectpolicy.DecodeConfig(pack)
}

func errMissingPackPolicy(packID string) (err error) {
	return fmt.Errorf("packs.%s policy is required", packID)
}
