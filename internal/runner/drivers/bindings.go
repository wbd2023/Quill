package drivers

import "ciphera/tools/internal/runner/drivers/internal/runtimebinding"

// RepositoryScanner is repository scanner.
type RepositoryScanner = runtimebinding.RepositoryScanner

// TargetCommand is target command.
type TargetCommand = runtimebinding.TargetCommand

// TargetCheck is target check.
type TargetCheck = runtimebinding.TargetCheck

// ProfileCheck is project check.
type ProfileCheck = runtimebinding.ProfileCheck

// Bindings is bindings.
type Bindings struct {
	repositoryScanners runtimebinding.RepositoryScanners
	targetCommands     runtimebinding.TargetCommands
	targetChecks       runtimebinding.TargetChecks
	projectChecks      runtimebinding.ProfileChecks
}

// NewBindings new bindings.
func NewBindings() (bindings Bindings) {
	return Bindings{
		repositoryScanners: runtimebinding.NewRepositoryScanners(),
		targetCommands:     runtimebinding.NewTargetCommands(),
		targetChecks:       runtimebinding.NewTargetChecks(),
		projectChecks:      runtimebinding.NewProjectChecks(),
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
