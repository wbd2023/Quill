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
// rule could not run). Output carries captured command output for fix failures.
type ExecutionResult struct {
	Diagnostics []Diagnostic

	ExitCode  int
	Output    string
	TimedOut  bool
	Truncated bool
}

// Empty reports whether the result has no diagnostics, command data, or output.
func (result ExecutionResult) Empty() (empty bool) {
	return len(result.Diagnostics) == 0 &&
		result.ExitCode == 0 &&
		!result.TimedOut &&
		!result.Truncated &&
		result.Output == ""
}

// HasCommand reports whether the result carries command execution metadata.
func (result ExecutionResult) HasCommand() (present bool) {
	return result.ExitCode != 0 || result.TimedOut || result.Truncated
}
