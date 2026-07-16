package style

// Diagnostic is a single style-check finding for one file and line.
type Diagnostic struct {
	Code    string
	File    string
	Line    int
	Column  int
	Message string
}

// ExecutionResult represents the outcome of running one check or fix against a rule.
type ExecutionResult struct {
	Diagnostics []Diagnostic

	ExitCode  int
	Output    string
	TimedOut  bool
	Truncated bool
}

// Empty reports whether the result has no diagnostics or command metadata.
func (result ExecutionResult) Empty() (empty bool) {
	return len(result.Diagnostics) == 0 &&
		result.ExitCode == 0 &&
		!result.TimedOut &&
		!result.Truncated
}

// HasCommand reports whether the result carries command execution metadata.
func (result ExecutionResult) HasCommand() (present bool) {
	return result.ExitCode != 0 || result.TimedOut || result.Truncated
}
