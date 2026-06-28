package style

// Diagnostic is a single style-check finding for one file and line.
type Diagnostic struct {
	Code    string
	File    string
	Line    int
	Column  int
	Message string
}

// ExecutionResult holds the outcome of running one check or fix against a rule: diagnostics, tool
// output, and the raw command result. A non-empty result signals violations; the error return is
// reserved for operational failures (the rule could not run).
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
