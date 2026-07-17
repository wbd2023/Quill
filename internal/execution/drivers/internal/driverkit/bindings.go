package driverkit

import (
	"context"
	"fmt"

	"ciphera/tools/internal/execution"
	"ciphera/tools/internal/process"
	"ciphera/tools/internal/style"
)

/* ----------------------------------------- Registries ----------------------------------------- */

// RepositoryScanner is repository scanner.
type RepositoryScanner func(
	ctx context.Context,
	run execution.RunContext,
	execution style.RepositoryScanExecution,
) (result style.ExecutionResult, err error)

// TargetCommand is target command.
type TargetCommand func(
	ctx context.Context,
	run execution.RunContext,
	job style.Job,
) (result style.ExecutionResult, err error)

// TargetCheck is target check.
type TargetCheck func(
	ctx context.Context,
	run execution.RunContext,
	job style.Job,
) (result style.ExecutionResult, err error)

// ProfileCheck is project check.
type ProfileCheck func(
	ctx context.Context,
	run execution.RunContext,
	execution style.ProfileExecution,
) (result style.ExecutionResult, err error)

// FileInterpreter converts a tool's raw command output into style diagnostics. It owns the
// tool-specific knowledge of which exit code signals findings and how the output is structured:
// the kind of knowledge that previously lived as FindingExitCode on FileCommandExecution. The
// interpreter returns nil diagnostics for a clean run. Non-nil error means the command genuinely
// failed (launch error, timeout, unexpected exit) and the diagnostics, if any, are diagnostic
// context for that failure.
type FileInterpreter func(result process.CommandResult) (diagnostics []style.Diagnostic, err error)

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

// ProfileChecks is project checks.
type ProfileChecks struct {
	entries map[string]ProfileCheck
}

// FileInterpreters is file interpreters.
type FileInterpreters struct {
	entries map[string]FileInterpreter
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
func NewProjectChecks() (registry ProfileChecks) {
	return ProfileChecks{entries: map[string]ProfileCheck{}}
}

// NewFileInterpreters new file interpreters.
func NewFileInterpreters() (registry FileInterpreters) {
	return FileInterpreters{entries: map[string]FileInterpreter{}}
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

func (registry *ProfileChecks) Add(id string, check ProfileCheck) {
	registry.entries = addBinding(registry.entries, "project check", id, check)
}

func (registry *FileInterpreters) Add(id string, interpreter FileInterpreter) {
	registry.entries = addBinding(registry.entries, "file interpreter", id, interpreter)
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

func (registry ProfileChecks) Lookup(id string) (check ProfileCheck, found bool) {
	check, found = registry.entries[id]
	return check, found
}

func (registry FileInterpreters) Lookup(id string) (interpreter FileInterpreter, found bool) {
	interpreter, found = registry.entries[id]
	return interpreter, found
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
