package drivers

import (
	projectdrivers "ciphera/tools/internal/runner/drivers/project"
	scandrivers "ciphera/tools/internal/runner/drivers/scan"
	targetdrivers "ciphera/tools/internal/runner/drivers/target"
)

/* -------------------------------------- Project Bindings -------------------------------------- */

// CheckProjectEnforcementLevels check project enforcement levels.
func CheckProjectEnforcementLevels() (check ProjectCheck) {
	return projectdrivers.CheckEnforcementLevels()
}

// CheckProjectExcludedDirectories check project excluded directories.
func CheckProjectExcludedDirectories() (check ProjectCheck) {
	return projectdrivers.CheckExcludedDirectories()
}

// CheckProjectCommands check project commands.
func CheckProjectCommands(projectPackID string) (check ProjectCheck) {
	return projectdrivers.CheckCommands(projectPackID)
}

/* -------------------------------------- Scanner Bindings -------------------------------------- */

// CheckGoArchitecture check go architecture.
func CheckGoArchitecture(goPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckGoArchitecture(goPackID)
}

// CheckASCII check a s c i i.
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

func RunGoFormat(
	goPackID string,
	goimportsToolID string,
	goLanguage string,
) (command TargetCommand) {
	return targetdrivers.RunGoFormat(goPackID, goimportsToolID, goLanguage)
}

func CheckGoStyle(goPackID string, goLanguage string) (check TargetCheck) {
	return targetdrivers.CheckGoStyle(goPackID, goLanguage)
}
