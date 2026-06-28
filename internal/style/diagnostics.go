package style

import "errors"

// Diagnostic is a single style-check finding for one file and line.
type Diagnostic struct {
	Code    string
	File    string
	Line    int
	Column  int
	Message string
}

// ExecutionResult holds the outcome of running one check or fix against a rule: diagnostics, tool
// output, and the raw command result.
type ExecutionResult struct {
	Diagnostics []Diagnostic
	Output      string
	Command     CommandResult
}

// Empty reports whether the result has no diagnostics, output, or command data.
func (result ExecutionResult) Empty() (empty bool) {
	return len(result.Diagnostics) == 0 &&
		result.Output == "" &&
		result.Command == CommandResult{}
}

// ViolationsFound is the sentinel error returned when a check produces at least one diagnostic.
func ViolationsFound() (err error) {
	return errors.New("violations found")
}
