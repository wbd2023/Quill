package contract

type CommandResult struct {
	ExitCode  int
	TimedOut  bool
	Truncated bool
}
