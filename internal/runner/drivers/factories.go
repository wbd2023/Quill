package drivers

import (
	projectdrivers "ciphera/tools/internal/runner/drivers/project"
	scandrivers "ciphera/tools/internal/runner/drivers/scan"
	targetdrivers "ciphera/tools/internal/runner/drivers/target"
)

func CheckProjectEnforcementLevels() (check ProjectCheck) {
	return projectdrivers.CheckEnforcementLevels()
}

func CheckProjectExcludedDirectories() (check ProjectCheck) {
	return projectdrivers.CheckExcludedDirectories()
}

func CheckProjectCommands(projectPackID string) (check ProjectCheck) {
	return projectdrivers.CheckCommands(projectPackID)
}

func CheckGoArchitecture(goPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckGoArchitecture(goPackID)
}

func CheckASCII() (scanner RepositoryScanner) {
	return scandrivers.CheckASCII()
}

func CheckExceptionMarkers() (scanner RepositoryScanner) {
	return scandrivers.CheckExceptionMarkers()
}

func CheckLineLengths() (scanner RepositoryScanner) {
	return scandrivers.CheckLineLengths()
}

func CheckMaintenanceMarkers() (scanner RepositoryScanner) {
	return scandrivers.CheckMaintenanceMarkers()
}

func CheckSectionHeaderNames(textPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckSectionHeaderNames(textPackID)
}

func CheckSectionHeaderDensity(textPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckSectionHeaderDensity(textPackID)
}

func CheckSectionHeaders(textPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckSectionHeaders(textPackID)
}

func CheckBashMagicValues() (scanner RepositoryScanner) {
	return scandrivers.CheckBashMagicValues()
}

func CheckBashSafety() (scanner RepositoryScanner) {
	return scandrivers.CheckBashSafety()
}

func CheckBashStructure() (scanner RepositoryScanner) {
	return scandrivers.CheckBashStructure()
}

func CheckBashTestHygiene() (scanner RepositoryScanner) {
	return scandrivers.CheckBashTestHygiene()
}

func CheckSecrets() (scanner RepositoryScanner) {
	return scandrivers.CheckSecrets()
}

func CheckVocabulary(vocabularyPackID string) (scanner RepositoryScanner) {
	return scandrivers.CheckVocabulary(vocabularyPackID)
}

func RunGolangci(
	goPackID string,
	golangciLintToolID string,
	goimportsToolID string,
	goLanguage string,
) (command TargetCommand) {
	return targetdrivers.RunGolangci(goPackID, golangciLintToolID, goimportsToolID, goLanguage)
}

func RunGoFormat(goPackID string, goimportsToolID string, goLanguage string) (command TargetCommand) {
	return targetdrivers.RunGoFormat(goPackID, goimportsToolID, goLanguage)
}

func CheckGoStyle(goPackID string, goLanguage string) (check TargetCheck) {
	return targetdrivers.CheckGoStyle(goPackID, goLanguage)
}
