package drivers

import (
	"github.com/wbd2023/Quill/internal/execution/drivers/command"
	profiledrivers "github.com/wbd2023/Quill/internal/execution/drivers/profile"
	scandrivers "github.com/wbd2023/Quill/internal/execution/drivers/scan"
	targetdrivers "github.com/wbd2023/Quill/internal/execution/drivers/target"
)

/* -------------------------------------- Profile Bindings -------------------------------------- */

// CheckProfileEnforcementLevels check project enforcement levels.
func CheckProfileEnforcementLevels() (check ProfileCheck) {
	return profiledrivers.CheckEnforcementLevels()
}

// CheckProfileExcludedDirectories check project excluded directories.
func CheckProfileExcludedDirectories() (check ProfileCheck) {
	return profiledrivers.CheckExcludedDirectories()
}

// CheckProfileCommands check project commands.
func CheckProfileCommands(profilePackID string) (check ProfileCheck) {
	return profiledrivers.CheckCommands(profilePackID)
}

/* -------------------------------------- Scanner Bindings -------------------------------------- */

// CheckGoArchitecture check go architecture.
func CheckGoArchitecture(goPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckGoArchitecture(goPackID)
}

// CheckASCII returns a scanner that flags non-ASCII characters.
func CheckASCII() (scanner RepositoryScanner) {
	return scandrivers.CheckASCII()
}

// CheckExceptionMarkers check exception markers.
func CheckExceptionMarkers() (scanner RepositoryScanner) {
	return scandrivers.CheckExceptionMarkers()
}

// CheckLineLengths check line lengths.
func CheckLineLengths() (scanner RepositoryScanner) {
	return scandrivers.CheckLineLengths()
}

// CheckMaintenanceMarkers check maintenance markers.
func CheckMaintenanceMarkers() (scanner RepositoryScanner) {
	return scandrivers.CheckMaintenanceMarkers()
}

// CheckSectionHeaderNames check section header names.
func CheckSectionHeaderNames(textPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckSectionHeaderNames(textPackID)
}

// CheckSectionHeaderDensity check section header density.
func CheckSectionHeaderDensity(textPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckSectionHeaderDensity(textPackID)
}

// CheckSectionHeaders check section headers.
func CheckSectionHeaders(textPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckSectionHeaders(textPackID)
}

// CheckBashMagicValues check bash magic values.
func CheckBashMagicValues() (scanner RepositoryScanner) {
	return scandrivers.CheckBashMagicValues()
}

// CheckBashSafety check bash safety.
func CheckBashSafety() (scanner RepositoryScanner) {
	return scandrivers.CheckBashSafety()
}

// CheckBashStructure check bash structure.
func CheckBashStructure() (scanner RepositoryScanner) {
	return scandrivers.CheckBashStructure()
}

// CheckBashTestHygiene check bash test hygiene.
func CheckBashTestHygiene() (scanner RepositoryScanner) {
	return scandrivers.CheckBashTestHygiene()
}

// CheckSecrets check secrets.
func CheckSecrets() (scanner RepositoryScanner) {
	return scandrivers.CheckSecrets()
}

// CheckVocabulary check vocabulary.
func CheckVocabulary(vocabularyPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckVocabulary(vocabularyPackID)
}

/* --------------------------------------- Target Bindings -------------------------------------- */

// RunGolangci run golangci.
func RunGolangci(
	goPackID string,
	golangciLintToolID string,
	goimportsToolID string,
	goLanguage string,
) (command TargetCommand) {
	return targetdrivers.RunGolangci(goPackID, golangciLintToolID, goimportsToolID, goLanguage)
}

// RunGoFormat go format.
func RunGoFormat(
	goPackID string,
	goimportsToolID string,
	goLanguage string,
) (command TargetCommand) {
	return targetdrivers.RunGoFormat(goPackID, goimportsToolID, goLanguage)
}

// CheckGoStyle go style.
func CheckGoStyle(goPackID string, goLanguage string) (check TargetCheck) {
	return targetdrivers.CheckGoStyle(goPackID, goLanguage)
}

/* ------------------------------------ Interpreter Bindings ------------------------------------ */

// ExitFindings is the conventional Unix linter findings exit code.
const ExitFindings = command.ExitFindings

// ExitFindingsMisspell is misspell's findings exit code when invoked with -error.
const ExitFindingsMisspell = command.ExitFindingsMisspell

// InterpretPlainText returns a file interpreter for tools whose findings output is multi-line
// text (shellcheck, markdownlint, misspell). When the tool exits with code, its trimmed output
// becomes a single diagnostic.
func InterpretPlainText(code int, codeLabel string) (interpreter FileInterpreter) {
	return command.InterpretPlainText(code, codeLabel)
}

// InterpretLines returns a file interpreter for tools whose findings output is one finding per
// line (gofmt -l, shfmt -d).
func InterpretLines(code int, codeLabel string) (interpreter FileInterpreter) {
	return command.InterpretLines(code, codeLabel)
}
