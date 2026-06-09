package drivers

import (
	"fmt"

	"ciphera/tools/internal/runner"
	commanddrivers "ciphera/tools/internal/runner/drivers/command"
	"ciphera/tools/internal/runner/drivers/internal/binding"
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

type RepositoryScanner = binding.RepositoryScanner

type TargetCommand = binding.TargetCommand

type TargetCheck = binding.TargetCheck

type ProjectCheck = binding.ProjectCheck

type Bindings struct {
	repositoryScanners binding.RepositoryScanners
	targetCommands     binding.TargetCommands
	targetChecks       binding.TargetChecks
	projectChecks      binding.ProjectChecks
}

func NewBindings() (bindings Bindings) {
	return Bindings{
		repositoryScanners: binding.NewRepositoryScanners(),
		targetCommands:     binding.NewTargetCommands(),
		targetChecks:       binding.NewTargetChecks(),
		projectChecks:      binding.NewProjectChecks(),
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
