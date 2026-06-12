package drivers

import "ciphera/tools/internal/runner/drivers/internal/runtimebinding"

type RepositoryScanner = runtimebinding.RepositoryScanner

type TargetCommand = runtimebinding.TargetCommand

type TargetCheck = runtimebinding.TargetCheck

type ProjectCheck = runtimebinding.ProjectCheck

type Bindings struct {
	repositoryScanners runtimebinding.RepositoryScanners
	targetCommands     runtimebinding.TargetCommands
	targetChecks       runtimebinding.TargetChecks
	projectChecks      runtimebinding.ProjectChecks
}

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

func (bindings *Bindings) AddProjectCheck(id string, check ProjectCheck) {
	bindings.projectChecks.Add(id, check)
}
