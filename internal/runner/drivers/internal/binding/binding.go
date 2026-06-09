package binding

import (
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

type RepositoryScanner func(
	context runner.Context,
	execution style.RepositoryScanExecution,
) (result style.ExecutionResult, err error)

type TargetCommand func(
	context runner.Context,
	spec style.ExecutionSpec,
) (result style.ExecutionResult, err error)

type TargetCheck func(
	context runner.Context,
	spec style.ExecutionSpec,
) (result style.ExecutionResult, err error)

type ProjectCheck func(
	context runner.Context,
	execution style.ProjectExecution,
) (result style.ExecutionResult, err error)

type RepositoryScanners struct {
	entries map[string]RepositoryScanner
}

type TargetCommands struct {
	entries map[string]TargetCommand
}

type TargetChecks struct {
	entries map[string]TargetCheck
}

type ProjectChecks struct {
	entries map[string]ProjectCheck
}

func NewRepositoryScanners() (registry RepositoryScanners) {
	return RepositoryScanners{entries: map[string]RepositoryScanner{}}
}

func NewTargetCommands() (registry TargetCommands) {
	return TargetCommands{entries: map[string]TargetCommand{}}
}

func NewTargetChecks() (registry TargetChecks) {
	return TargetChecks{entries: map[string]TargetCheck{}}
}

func NewProjectChecks() (registry ProjectChecks) {
	return ProjectChecks{entries: map[string]ProjectCheck{}}
}

func (registry *RepositoryScanners) Add(id string, scanner RepositoryScanner) {
	registry.entries = addBinding(registry.entries, "repository scanner", id, scanner)
}

func (registry *TargetCommands) Add(id string, command TargetCommand) {
	registry.entries = addBinding(registry.entries, "target command", id, command)
}

func (registry *TargetChecks) Add(id string, check TargetCheck) {
	registry.entries = addBinding(registry.entries, "target check", id, check)
}

func (registry *ProjectChecks) Add(id string, check ProjectCheck) {
	registry.entries = addBinding(registry.entries, "project check", id, check)
}

func (registry RepositoryScanners) Lookup(id string) (scanner RepositoryScanner, found bool) {
	scanner, found = registry.entries[id]
	return scanner, found
}

func (registry TargetCommands) Lookup(id string) (command TargetCommand, found bool) {
	command, found = registry.entries[id]
	return command, found
}

func (registry TargetChecks) Lookup(id string) (check TargetCheck, found bool) {
	check, found = registry.entries[id]
	return check, found
}

func (registry ProjectChecks) Lookup(id string) (check ProjectCheck, found bool) {
	check, found = registry.entries[id]
	return check, found
}

func addBinding[T any](entries map[string]T, kind string, id string, entry T) (updated map[string]T) {
	if entries == nil {
		entries = map[string]T{}
	}

	if _, exists := entries[id]; exists {
		panic(fmt.Sprintf("duplicate %s binding %q", kind, id))
	}

	entries[id] = entry
	return entries
}
