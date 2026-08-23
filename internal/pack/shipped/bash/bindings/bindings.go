// Package bindings owns the Bash Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Bash execution identities (scanners and file interpreters)
// to concrete check behaviour. The parent bash package stays independent of the driver facade and
// check implementations.
package bindings

import (
	"context"

	checks "github.com/wbd2023/quill/internal/checks/bash"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/bash"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/style"
)

// Register wires every Bash execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	registerFileInterpreters(bindings)
	registerRepositoryScanners(bindings)
}

func registerFileInterpreters(bindings *drivers.Bindings) {
	bindings.AddFileInterpreter(
		tool.Shellcheck,
		drivers.InterpretPlainText(drivers.ExitFindings, "bash/shellcheck/findings"),
	)
	bindings.AddFileInterpreter(
		tool.Shfmt,
		drivers.InterpretLines(drivers.ExitFindings, "bash/shfmt/findings"),
	)
}

func registerRepositoryScanners(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(bash.PackID, bash.ScannerMagicValues, scanMagicValues)
	bindings.AddRepositoryScanner(bash.PackID, bash.ScannerSafety, scanSafety)
	bindings.AddRepositoryScanner(bash.PackID, bash.ScannerStructure, scanStructure)
	bindings.AddRepositoryScanner(bash.PackID, bash.ScannerTestHygiene, scanTestHygiene)
}

func scanMagicValues(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckMagicValues(
		context.RepoRoot,
		context.Profile.Repository,
		context.Scope,
	)
}

func scanSafety(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckSafety(context.RepoRoot, context.Profile.Repository, context.Scope)
}

func scanStructure(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckStructure(
		context.RepoRoot,
		context.Profile.Repository,
		context.Scope,
	)
}

func scanTestHygiene(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckTestHygiene(
		context.RepoRoot,
		context.Profile.Repository,
		context.Scope,
	)
}
