package style

// Diagnostic is a single style-check finding located at a Range within a repository-relative File.
// The zero Range represents an unknown location; see VerifyRange for the protocol-boundary check
// applied to diagnostics entering Quill from external sources. HelpURL carries an optional
// documentation link supplied by the producer; it is empty when none is provided.
type Diagnostic struct {
	Code    string
	Message string

	File    string
	Range   Range
	HelpURL string
}

// ExecutionResult represents the outcome of running one check or fix against a rule. Pack
// provenance is carried by the Rule that produced the result, not duplicated here: report
// aggregation attributes findings through the rule association.
type ExecutionResult struct {
	Diagnostics []Diagnostic

	ExitCode int
	Output   string

	TimedOut  bool
	Truncated bool
}

// Empty reports whether the result has no diagnostics or command metadata.
func (result ExecutionResult) Empty() (empty bool) {
	return len(result.Diagnostics) == 0 &&
		result.ExitCode == 0 &&
		!result.TimedOut && !result.Truncated
}

// HasCommand reports whether the result carries command execution metadata.
func (result ExecutionResult) HasCommand() (present bool) {
	return result.ExitCode != 0 || result.TimedOut || result.Truncated
}
