package style

// Diagnostic is a single style-check finding for one file and line.
type Diagnostic struct {
	Code    string
	File    string
	Line    int
	Column  int
	Message string
}

// ExecutionResult holds the outcome of running one check or fix against a rule. A non-empty
// Diagnostics slice signals findings; the error return is reserved for operational failures (the
// rule could not run). There is no text-blob field: all findings are structured diagnostics.
type ExecutionResult struct {
	Diagnostics []Diagnostic
	Command     CommandResult
}

// Empty reports whether the result has no diagnostics or command data.
func (result ExecutionResult) Empty() (empty bool) {
	return len(result.Diagnostics) == 0 &&
		result.Command == CommandResult{}
}
