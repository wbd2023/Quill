package drivers

import "ciphera/tools/internal/execution/drivers/internal/runtimebinding"

// RepositoryScanner is repository scanner.
type RepositoryScanner = runtimebinding.RepositoryScanner

// TargetCommand is target command.
type TargetCommand = runtimebinding.TargetCommand

// TargetCheck is target check.
type TargetCheck = runtimebinding.TargetCheck

// ProfileCheck is project check.
type ProfileCheck = runtimebinding.ProfileCheck

// FileInterpreter converts a tool's raw command output into style diagnostics.
type FileInterpreter = runtimebinding.FileInterpreter

// Bindings is bindings.
type Bindings struct {
	repositoryScanners runtimebinding.RepositoryScanners
	targetCommands     runtimebinding.TargetCommands
	targetChecks       runtimebinding.TargetChecks
	projectChecks      runtimebinding.ProfileChecks
	fileInterpreters   runtimebinding.FileInterpreters
}

// NewBindings new bindings.
func NewBindings() (bindings Bindings) {
	return Bindings{
		repositoryScanners: runtimebinding.NewRepositoryScanners(),
		targetCommands:     runtimebinding.NewTargetCommands(),
		targetChecks:       runtimebinding.NewTargetChecks(),
		projectChecks:      runtimebinding.NewProjectChecks(),
		fileInterpreters:   runtimebinding.NewFileInterpreters(),
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
