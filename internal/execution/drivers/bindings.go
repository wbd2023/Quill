package drivers

import "ciphera/tools/internal/execution/drivers/internal/driverkit"

// RepositoryScanner is repository scanner.
type RepositoryScanner = driverkit.RepositoryScanner

// TargetCommand is target command.
type TargetCommand = driverkit.TargetCommand

// TargetCheck is target check.
type TargetCheck = driverkit.TargetCheck

// ProfileCheck is project check.
type ProfileCheck = driverkit.ProfileCheck

// FileInterpreter converts a tool's raw command output into style diagnostics.
type FileInterpreter = driverkit.FileInterpreter

// Bindings is bindings.
type Bindings struct {
	repositoryScanners driverkit.RepositoryScanners
	targetCommands     driverkit.TargetCommands
	targetChecks       driverkit.TargetChecks
	projectChecks      driverkit.ProfileChecks
	fileInterpreters   driverkit.FileInterpreters
}

// NewBindings new bindings.
func NewBindings() (bindings Bindings) {
	return Bindings{
		repositoryScanners: driverkit.NewRepositoryScanners(),
		targetCommands:     driverkit.NewTargetCommands(),
		targetChecks:       driverkit.NewTargetChecks(),
		projectChecks:      driverkit.NewProjectChecks(),
		fileInterpreters:   driverkit.NewFileInterpreters(),
	}
}

func (bindings *Bindings) AddRepositoryScanner(id string, scanner RepositoryScanner) {
	bindings.repositoryScanners.Add(id, scanner)
}

func (bindings *Bindings) AddTargetCommand(action string, command TargetCommand) {
	bindings.targetCommands.Add(action, command)
}

func (bindings *Bindings) AddTargetCheck(language string, check TargetCheck) {
	bindings.targetChecks.Add(language, check)
}

func (bindings *Bindings) AddProjectCheck(id string, check ProfileCheck) {
	bindings.projectChecks.Add(id, check)
}

func (bindings *Bindings) AddFileInterpreter(id string, interpreter FileInterpreter) {
	bindings.fileInterpreters.Add(id, interpreter)
}

func (bindings Bindings) LookupFileInterpreter(id string) (
	interpreter FileInterpreter,
	found bool,
) {
	return bindings.fileInterpreters.Lookup(id)
}
