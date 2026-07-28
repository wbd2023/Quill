package drivers

import "github.com/wbd2023/quill/internal/style"

// Bindings is the aggregate of every shipped runtime binding. Repository scanners, Profile checks,
// target commands, and target checks are keyed by Pack-qualified identity; file interpreters stay
// keyed by global Tool ID. Bindings satisfies style.RuntimeBindings so the completeness validation
// can verify every active execution resolves to exactly one registered binding.
type Bindings struct {
	repositoryScanners RepositoryScanners
	targetCommands     TargetCommands
	targetChecks       TargetChecks
	profileChecks      ProfileChecks
	fileInterpreters   FileInterpreters
}

// Compile-time assertion that Bindings satisfies the Pack-qualified runtime binding contract.
var _ style.RuntimeBindings = Bindings{}

// NewBindings new bindings.
func NewBindings() (bindings Bindings) {
	return Bindings{
		repositoryScanners: NewRepositoryScanners(),
		targetCommands:     NewTargetCommands(),
		targetChecks:       NewTargetChecks(),
		profileChecks:      NewProfileChecks(),
		fileInterpreters:   NewFileInterpreters(),
	}
}

func (bindings *Bindings) AddRepositoryScanner(
	packID string,
	id string,
	scanner RepositoryScanner,
) {
	bindings.repositoryScanners.Add(packID, id, scanner)
}

func (bindings *Bindings) AddTargetCommand(
	packID string,
	language string,
	action string,
	command TargetCommand,
) {
	bindings.targetCommands.Add(packID, language, action, command)
}

func (bindings *Bindings) AddTargetCheck(
	packID string,
	language string,
	id string,
	check TargetCheck,
) {
	bindings.targetChecks.Add(packID, language, id, check)
}

func (bindings *Bindings) AddProfileCheck(packID string, id string, check ProfileCheck) {
	bindings.profileChecks.Add(packID, id, check)
}

func (bindings *Bindings) AddFileInterpreter(id string, interpreter FileInterpreter) {
	bindings.fileInterpreters.Add(id, interpreter)
}

func (bindings Bindings) LookupRepositoryScanner(packID string, id string) (
	scanner RepositoryScanner,
	found bool,
) {
	return bindings.repositoryScanners.Lookup(packID, id)
}

func (bindings Bindings) LookupTargetCommand(
	packID string,
	language string,
	action string,
) (command TargetCommand, found bool) {
	return bindings.targetCommands.Lookup(packID, language, action)
}

func (bindings Bindings) LookupTargetCheck(
	packID string,
	language string,
	id string,
) (check TargetCheck, found bool) {
	return bindings.targetChecks.Lookup(packID, language, id)
}

func (bindings Bindings) LookupProfileCheck(packID string, id string) (
	check ProfileCheck,
	found bool,
) {
	return bindings.profileChecks.Lookup(packID, id)
}

func (bindings Bindings) LookupFileInterpreter(id string) (
	interpreter FileInterpreter,
	found bool,
) {
	return bindings.fileInterpreters.Lookup(id)
}

/* ---------------------------------- Runtime Bindings Contract --------------------------------- */

// HasProfileCheck reports whether a Profile check is registered for the Pack-qualified identity.
func (bindings Bindings) HasProfileCheck(packID string, check string) (found bool) {
	_, found = bindings.profileChecks.Lookup(packID, check)
	return found
}

// HasRepositoryScanner reports whether a repository scanner is registered for the Pack-qualified
// identity.
func (bindings Bindings) HasRepositoryScanner(packID string, scanner string) (found bool) {
	_, found = bindings.repositoryScanners.Lookup(packID, scanner)
	return found
}

// HasTargetCommand reports whether a target command is registered for the Pack-qualified identity.
func (bindings Bindings) HasTargetCommand(
	packID string,
	language string,
	action string,
) (found bool) {
	_, found = bindings.targetCommands.Lookup(packID, language, action)
	return found
}

// HasTargetCheck reports whether a target check is registered for the Pack-qualified identity.
func (bindings Bindings) HasTargetCheck(packID string, language string, check string) (found bool) {
	_, found = bindings.targetChecks.Lookup(packID, language, check)
	return found
}

// HasFileInterpreter reports whether a file interpreter is registered for the global Tool ID.
func (bindings Bindings) HasFileInterpreter(toolID string) (found bool) {
	_, found = bindings.fileInterpreters.Lookup(toolID)
	return found
}
