package bindings

import (
	"github.com/wbd2023/Quill/internal/execution/drivers"
	"github.com/wbd2023/Quill/internal/pack/shipped/bash"
	"github.com/wbd2023/Quill/internal/pack/shipped/golang"
	"github.com/wbd2023/Quill/internal/pack/shipped/project"
	"github.com/wbd2023/Quill/internal/pack/shipped/security"
	"github.com/wbd2023/Quill/internal/pack/shipped/text"
	"github.com/wbd2023/Quill/internal/pack/shipped/tool"
	"github.com/wbd2023/Quill/internal/pack/shipped/vocabulary"
)

// Build wires every shipped pack's scanners, commands, checks, and interpreters into a single
// Bindings value for driver construction.
func Build() (bindings drivers.Bindings) {
	bindings = drivers.NewBindings()
	registerFileInterpreters(&bindings)
	registerProfileChecks(&bindings)
	registerRepositoryScanners(&bindings)
	registerTargetBindings(&bindings)
	return bindings
}

/* -------------------------------------- File Interpreters ------------------------------------- */

func registerFileInterpreters(bindings *drivers.Bindings) {
	bindings.AddFileInterpreter(
		tool.Shellcheck,
		drivers.InterpretPlainText(drivers.ExitFindings, "bash/shellcheck/findings"),
	)
	bindings.AddFileInterpreter(
		tool.Shfmt,
		drivers.InterpretLines(drivers.ExitFindings, "bash/shfmt/findings"),
	)
	bindings.AddFileInterpreter(
		tool.Misspell,
		drivers.InterpretPlainText(drivers.ExitFindingsMisspell, "text/spelling/findings"),
	)
	bindings.AddFileInterpreter(
		tool.Markdownlint,
		drivers.InterpretPlainText(drivers.ExitFindings, "markdown/markdownlint/findings"),
	)
}

/* --------------------------------------- Profile Checks --------------------------------------- */

func registerProfileChecks(bindings *drivers.Bindings) {
	bindings.AddProfileCheck(
		project.CheckEnforcementLevels,
		drivers.CheckProfileEnforcementLevels(),
	)
	bindings.AddProfileCheck(
		project.CheckExcludedDirectories,
		drivers.CheckProfileExcludedDirectories(),
	)
	bindings.AddProfileCheck(
		project.CheckCommands,
		drivers.CheckProfileCommands(project.PackID),
	)
}

/* ------------------------------------- Repository Scanners ------------------------------------ */

func registerRepositoryScanners(bindings *drivers.Bindings) {
	bindings.AddRepositoryScanner(
		golang.ScannerArchitecture,
		drivers.CheckGoArchitecture(golang.PackID),
	)
	bindings.AddRepositoryScanner(text.ScannerASCII, drivers.CheckASCII())
	bindings.AddRepositoryScanner(
		text.ScannerExceptionMarkers,
		drivers.CheckExceptionMarkers(),
	)
	bindings.AddRepositoryScanner(text.ScannerLineLength, drivers.CheckLineLengths())
	bindings.AddRepositoryScanner(
		text.ScannerMaintenanceMarkers,
		drivers.CheckMaintenanceMarkers(),
	)
	bindings.AddRepositoryScanner(
		text.ScannerSectionHeaderNames,
		drivers.CheckSectionHeaderNames(text.PackID),
	)
	bindings.AddRepositoryScanner(
		text.ScannerSectionHeaderDensity,
		drivers.CheckSectionHeaderDensity(text.PackID),
	)
	bindings.AddRepositoryScanner(
		text.ScannerSectionHeaders,
		drivers.CheckSectionHeaders(text.PackID),
	)
	bindings.AddRepositoryScanner(
		bash.ScannerMagicValues,
		drivers.CheckBashMagicValues(),
	)
	bindings.AddRepositoryScanner(bash.ScannerSafety, drivers.CheckBashSafety())
	bindings.AddRepositoryScanner(bash.ScannerStructure, drivers.CheckBashStructure())
	bindings.AddRepositoryScanner(
		bash.ScannerTestHygiene,
		drivers.CheckBashTestHygiene(),
	)
	bindings.AddRepositoryScanner(security.ScannerSecrets, drivers.CheckSecrets())
	bindings.AddRepositoryScanner(
		vocabulary.ScannerVocabulary,
		drivers.CheckVocabulary(vocabulary.PackID),
	)
}

/* --------------------------------------- Target Bindings -------------------------------------- */

func registerTargetBindings(bindings *drivers.Bindings) {
	bindings.AddTargetCommand(
		golang.TargetActionGolangci,
		drivers.RunGolangci(
			golang.PackID,
			tool.GolangciLint,
			tool.Goimports,
			golang.Language,
		),
	)
	bindings.AddTargetCommand(
		golang.TargetActionGoFormat,
		drivers.RunGoFormat(golang.PackID, tool.Goimports, golang.Language),
	)
	bindings.AddTargetCheck(golang.Language, drivers.CheckGoStyle(golang.PackID, golang.Language))
}
