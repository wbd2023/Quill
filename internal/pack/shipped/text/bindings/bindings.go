// Package bindings owns the Text Shipped Pack's runtime driver registrations.
//
// It is the only place that may connect Text execution identities (the misspell file interpreter
// and the text repository scanners) to concrete check behaviour. The parent text package stays
// independent of the driver facade and check implementations.
package bindings

import (
	"context"
	"fmt"

	checks "github.com/wbd2023/quill/internal/checks/text"
	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/execution/drivers"
	"github.com/wbd2023/quill/internal/pack/shipped/text"
	textpolicy "github.com/wbd2023/quill/internal/pack/shipped/text/policy"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/style"
)

/* ---------------------------------------- Registration ---------------------------------------- */

// Register wires every Text execution identity into the aggregate driver Bindings.
// It is called only by the central shipped aggregate builder.
func Register(bindings *drivers.Bindings) {
	registerFileInterpreters(bindings)
	registerRepositoryScanners(bindings)
}

func registerFileInterpreters(bindings *drivers.Bindings) {
	bindings.AddFileInterpreter(
		tool.Misspell,
		drivers.InterpretPlainText(drivers.ExitFindingsMisspell, "text/spelling/findings"),
	)
}

func registerRepositoryScanners(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(text.PackID, text.ScannerASCII, scanASCII)
	bindings.AddRepositoryScanner(text.PackID, text.ScannerExceptionMarkers, scanExceptionMarkers)
	bindings.AddRepositoryScanner(text.PackID, text.ScannerLineLength, scanLineLengths)
	bindings.AddRepositoryScanner(
		text.PackID,
		text.ScannerMaintenanceMarkers,
		scanMaintenanceMarkers,
	)
	bindings.AddRepositoryScanner(
		text.PackID,
		text.ScannerSectionHeaderNames,
		scanSectionHeaderNames,
	)
	bindings.AddRepositoryScanner(
		text.PackID,
		text.ScannerSectionHeaderDensity,
		scanSectionHeaderDensity,
	)
	bindings.AddRepositoryScanner(text.PackID, text.ScannerSectionHeaders, scanSectionHeaders)
}

/* ------------------------------------- Repository Scanners ------------------------------------ */

func scanASCII(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckASCII(context.Root, context.Profile.Repository, context.Scope)
}

func scanExceptionMarkers(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckExceptionMarkers(
		context.Root,
		context.Profile.Repository,
		context.Scope,
	)
}

func scanLineLengths(
	_ context.Context,
	context execution.RunContext,
	scanExec style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	files, err := execution.CollectFileSetFiles(context, scanExec.FileSet)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return checks.CheckLineLengths(context.Root, files)
}

func scanMaintenanceMarkers(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	return checks.CheckMaintenanceMarkers(
		context.Root,
		context.Profile.Repository,
		context.Scope,
	)
}

func scanSectionHeaderNames(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	config, err := decodeTextPackPolicy(context, text.PackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return checks.CheckSectionHeaderNames(
		context.Root,
		context.Profile.Repository,
		config.SectionHeaders,
		context.Scope,
	)
}

func scanSectionHeaderDensity(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	config, err := decodeTextPackPolicy(context, text.PackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return checks.CheckSectionHeaderDensity(
		context.Root,
		context.Profile.Repository,
		config.SectionHeaders,
		context.Scope,
	)
}

func scanSectionHeaders(
	_ context.Context,
	context execution.RunContext,
	_ style.RepositoryScan,
) (result style.ExecutionResult, err error) {
	config, err := decodeTextPackPolicy(context, text.PackID)
	if err != nil {
		return style.ExecutionResult{}, err
	}

	return checks.CheckSectionHeaders(
		context.Root,
		context.Profile.Repository,
		config.SectionHeaders,
		context.Scope,
	)
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func decodeTextPackPolicy(
	context execution.RunContext,
	packID string,
) (config textpolicy.Config, err error) {
	pack, found := context.Profile.PackPolicies.Lookup(packID)
	if !found {
		return textpolicy.Config{}, errMissingPackPolicy(packID)
	}

	return textpolicy.DecodeConfig(pack)
}

func errMissingPackPolicy(packID string) (err error) {
	return fmt.Errorf("packs.%s policy is required", packID)
}
