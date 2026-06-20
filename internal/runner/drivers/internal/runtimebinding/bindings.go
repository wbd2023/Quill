package runtimebinding

import (
	"fmt"

	"ciphera/tools/internal/runner"
	"ciphera/tools/internal/style"
)

/* ----------------------------------------- Registries ----------------------------------------- */

// RepositoryScanner is repository scanner.
type RepositoryScanner func(
	context runner.Context,
	execution style.RepositoryScanExecution,
) (result style.ExecutionResult, err error)

// TargetCommand is target command.
type TargetCommand func(
	context runner.Context,
	spec style.ExecutionSpec,
) (result style.ExecutionResult, err error)

// TargetCheck is target check.
type TargetCheck func(
	context runner.Context,
	spec style.ExecutionSpec,
) (result style.ExecutionResult, err error)

// ProjectCheck is project check.
type ProjectCheck func(
	context runner.Context,
	execution style.ProjectExecution,
) (result style.ExecutionResult, err error)

// RepositoryScanners is repository scanners.
type RepositoryScanners struct {
	entries map[string]RepositoryScanner
}

// TargetCommands is target commands.
type TargetCommands struct {
	entries map[string]TargetCommand
}

// TargetChecks is target checks.
type TargetChecks struct {
	entries map[string]TargetCheck
}

// ProjectChecks is project checks.
type ProjectChecks struct {
	entries map[string]ProjectCheck
}

// NewRepositoryScanners new repository scanners.
func NewRepositoryScanners() (registry RepositoryScanners) {
	return RepositoryScanners{entries: map[string]RepositoryScanner{}}
}

// NewTargetCommands new target commands.
func NewTargetCommands() (registry TargetCommands) {
	return TargetCommands{entries: map[string]TargetCommand{}}
}

// NewTargetChecks new target checks.
func NewTargetChecks() (registry TargetChecks) {
	return TargetChecks{entries: map[string]TargetCheck{}}
}

// NewProjectChecks new project checks.
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

/* ------------------------------------------- Helpers ------------------------------------------ */

func addBinding[T any](
	entries map[string]T,
	kind string,
	id string,
	entry T,
) (updated map[string]T) {
	if entries == nil {
		entries = map[string]T{}
	}

	if _, exists := entries[id]; exists {
		panic(fmt.Sprintf("duplicate %s binding %q", kind, id))
	}

	entries[id] = entry
	return entries
}
