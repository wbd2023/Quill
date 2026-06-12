package drivers

import (
	"fmt"

	"ciphera/tools/internal/runner"
	commanddrivers "ciphera/tools/internal/runner/drivers/command"
	"ciphera/tools/internal/runner/drivers/internal/runtimebinding"
	projectdrivers "ciphera/tools/internal/runner/drivers/project"
	scandrivers "ciphera/tools/internal/runner/drivers/scan"
	targetdrivers "ciphera/tools/internal/runner/drivers/target"
	"ciphera/tools/internal/style"
)

func CheckDrivers(bindings Bindings) (registry runner.DriverRegistry) {
	return mergeDrivers(
		runner.DriverRegistry{
			style.ExecutionToolchain: runner.ToolchainDriver,
		},
		commanddrivers.CheckDrivers(),
		projectdrivers.CheckDrivers(bindings.projectChecks),
		targetdrivers.CheckDrivers(bindings.targetCommands, bindings.targetChecks),
		scandrivers.CheckDrivers(bindings.repositoryScanners),
	)
}

func FixDrivers(bindings Bindings) (registry runner.DriverRegistry) {
	return mergeDrivers(
		commanddrivers.FixDrivers(),
		targetdrivers.FixDrivers(bindings.targetCommands),
	)
}

func mergeDrivers(registries ...runner.DriverRegistry) (merged runner.DriverRegistry) {
	merged = runner.DriverRegistry{}
	for _, registry := range registries {
		for kind, driver := range registry {
			if _, exists := merged[kind]; exists {
				panic(fmt.Sprintf("duplicate driver for execution kind %q", kind))
			}

			merged[kind] = driver
		}
	}

	return merged
}

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
