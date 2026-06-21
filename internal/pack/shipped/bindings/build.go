package bindings

import (
	"ciphera/tools/internal/pack/shipped/bash"
	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/project"
	"ciphera/tools/internal/pack/shipped/security"
	"ciphera/tools/internal/pack/shipped/text"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/pack/shipped/vocabulary"
	"ciphera/tools/internal/runner/drivers"
)

// Build returns the requested value.
func Build() (bindings drivers.Bindings) {
	bindings = drivers.NewBindings()

	bindings.AddProjectCheck(
		project.CheckEnforcementLevels,
		drivers.CheckProfileEnforcementLevels(),
	)
	bindings.AddProjectCheck(
		project.CheckExcludedDirectories,
		drivers.CheckProfileExcludedDirectories(),
	)
	bindings.AddProjectCheck(
		project.CheckCommands,
		drivers.CheckProfileCommands(project.PackID),
	)

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

	return bindings
}
