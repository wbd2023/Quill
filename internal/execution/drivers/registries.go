package drivers

import (
	"context"
	"fmt"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/process"
	"github.com/wbd2023/quill/internal/style"
)

/* ----------------------------------------- Registries ----------------------------------------- */

// RepositoryScanner runs one repository scan for a run context and returns its diagnostics as an
// ExecutionResult.
type RepositoryScanner func(
	ctx context.Context,
	run execution.RunContext,
	scan style.RepositoryScan,
) (result style.ExecutionResult, err error)

// TargetCommand is a bound target command. The resolved target names are carried by the Job, not
// declared by the Pack, so the closure receives one self-contained typed value.
type TargetCommand func(
	ctx context.Context,
	run execution.RunContext,
	job style.TargetCommandJob,
) (result style.ExecutionResult, err error)

// TargetCheck is a bound target check. The resolved target names are carried by the Job, not
// declared by the Pack, so the closure receives one self-contained typed value.
type TargetCheck func(
	ctx context.Context,
	run execution.RunContext,
	job style.TargetCheckJob,
) (result style.ExecutionResult, err error)

// ProfileCheck runs one Profile check for a run context and returns its diagnostics as an
// ExecutionResult.
type ProfileCheck func(
	ctx context.Context,
	run execution.RunContext,
	check style.ProfileCheck,
) (result style.ExecutionResult, err error)

// FileInterpreter converts a tool's raw command output into style diagnostics. It owns the
// tool-specific knowledge of which exit code signals findings and how the output is structured.
// File interpreters stay keyed by global Tool ID because a tool's output format is fixed by the
// tool, not by the Pack that declares a file-command rule.
type FileInterpreter func(result process.CommandResult) (diagnostics []style.Diagnostic, err error)

// scannerKey is the Pack-qualified identity of a repository scanner.
type scannerKey struct {
	packID string
	id     string
}

// checkKey is the Pack-qualified identity of a Profile check.
type checkKey struct {
	packID string
	id     string
}

// targetKey is the Pack-qualified identity of a target command or check.
type targetKey struct {
	packID   string
	language string
	id       string
}

// RepositoryScanners is repository scanners keyed by Pack-qualified scanner identity.
type RepositoryScanners struct {
	entries map[scannerKey]RepositoryScanner
}

// TargetCommands is target commands keyed by Pack-qualified (language, action) identity.
type TargetCommands struct {
	entries map[targetKey]TargetCommand
}

// TargetChecks is target checks keyed by Pack-qualified (language, local id) identity.
type TargetChecks struct {
	entries map[targetKey]TargetCheck
}

// ProfileChecks contains Profile checks keyed by Pack-qualified check identity.
type ProfileChecks struct {
	entries map[checkKey]ProfileCheck
}

// FileInterpreters is file interpreters keyed by global Tool ID.
type FileInterpreters struct {
	entries map[string]FileInterpreter
}

// NewRepositoryScanners returns an empty repository scanner registry keyed by Pack-qualified
// scanner identity.
func NewRepositoryScanners() (registry RepositoryScanners) {
	return RepositoryScanners{entries: map[scannerKey]RepositoryScanner{}}
}

// NewTargetCommands returns an empty target command registry keyed by Pack-qualified
// (language, action) identity.
func NewTargetCommands() (registry TargetCommands) {
	return TargetCommands{entries: map[targetKey]TargetCommand{}}
}

// NewTargetChecks returns an empty target check registry keyed by Pack-qualified
// (language, local id) identity.
func NewTargetChecks() (registry TargetChecks) {
	return TargetChecks{entries: map[targetKey]TargetCheck{}}
}

// NewProfileChecks returns an empty Profile Check registry.
func NewProfileChecks() (registry ProfileChecks) {
	return ProfileChecks{entries: map[checkKey]ProfileCheck{}}
}

// NewFileInterpreters returns an empty file interpreter registry keyed by global Tool ID.
func NewFileInterpreters() (registry FileInterpreters) {
	return FileInterpreters{entries: map[string]FileInterpreter{}}
}

func (registry *RepositoryScanners) Add(packID string, id string, scanner RepositoryScanner) {
	key := scannerKey{packID: packID, id: id}
	registry.entries = addScannerBinding(registry.entries, "repository scanner", key, scanner)
}

func (registry *TargetCommands) Add(
	packID string,
	language string,
	action string,
	command TargetCommand,
) {
	key := targetKey{packID: packID, language: language, id: action}
	registry.entries = addTargetBinding(registry.entries, "target command", key, command)
}

func (registry *TargetChecks) Add(
	packID string,
	language string,
	id string,
	check TargetCheck,
) {
	key := targetKey{packID: packID, language: language, id: id}
	registry.entries = addTargetBinding(registry.entries, "target check", key, check)
}

func (registry *ProfileChecks) Add(packID string, id string, check ProfileCheck) {
	key := checkKey{packID: packID, id: id}
	registry.entries = addCheckBinding(registry.entries, "profile check", key, check)
}

func (registry *FileInterpreters) Add(id string, interpreter FileInterpreter) {
	registry.entries = addInterpreterBinding(registry.entries, "file interpreter", id, interpreter)
}

func (registry RepositoryScanners) Lookup(packID string, id string) (
	scanner RepositoryScanner,
	found bool,
) {
	scanner, found = registry.entries[scannerKey{packID: packID, id: id}]
	return scanner, found
}

func (registry TargetCommands) Lookup(
	packID string,
	language string,
	action string,
) (command TargetCommand, found bool) {
	command, found = registry.entries[targetKey{packID: packID, language: language, id: action}]
	return command, found
}

func (registry TargetChecks) Lookup(
	packID string,
	language string,
	id string,
) (check TargetCheck, found bool) {
	check, found = registry.entries[targetKey{packID: packID, language: language, id: id}]
	return check, found
}

func (registry ProfileChecks) Lookup(packID string, id string) (check ProfileCheck, found bool) {
	check, found = registry.entries[checkKey{packID: packID, id: id}]
	return check, found
}

func (registry FileInterpreters) Lookup(id string) (interpreter FileInterpreter, found bool) {
	interpreter, found = registry.entries[id]
	return interpreter, found
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func addScannerBinding(
	entries map[scannerKey]RepositoryScanner,
	kind string,
	key scannerKey,
	entry RepositoryScanner,
) (updated map[scannerKey]RepositoryScanner) {
	if entries == nil {
		entries = map[scannerKey]RepositoryScanner{}
	}

	if _, exists := entries[key]; exists {
		panic(fmt.Sprintf("duplicate %s binding %q/%q", kind, key.packID, key.id))
	}

	entries[key] = entry
	return entries
}

func addTargetBinding[T any](
	entries map[targetKey]T,
	kind string,
	key targetKey,
	entry T,
) (updated map[targetKey]T) {
	if entries == nil {
		entries = map[targetKey]T{}
	}

	if _, exists := entries[key]; exists {
		panic(fmt.Sprintf("duplicate %s binding %q/%q/%q", kind, key.packID, key.language, key.id))
	}

	entries[key] = entry
	return entries
}

func addCheckBinding(
	entries map[checkKey]ProfileCheck,
	kind string,
	key checkKey,
	entry ProfileCheck,
) (updated map[checkKey]ProfileCheck) {
	if entries == nil {
		entries = map[checkKey]ProfileCheck{}
	}

	if _, exists := entries[key]; exists {
		panic(fmt.Sprintf("duplicate %s binding %q/%q", kind, key.packID, key.id))
	}

	entries[key] = entry
	return entries
}

func addInterpreterBinding(
	entries map[string]FileInterpreter,
	kind string,
	id string,
	entry FileInterpreter,
) (updated map[string]FileInterpreter) {
	if entries == nil {
		entries = map[string]FileInterpreter{}
	}

	if _, exists := entries[id]; exists {
		panic(fmt.Sprintf("duplicate %s binding %q", kind, id))
	}

	entries[id] = entry
	return entries
}
