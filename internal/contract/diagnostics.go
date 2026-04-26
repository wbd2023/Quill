package contract

import "errors"

type Diagnostic struct {
	Code    string
	File    string
	Line    int
	Column  int
	Message string
}

type ExecutionResult struct {
	Diagnostics []Diagnostic
	Output      string
	Command     CommandResult
}

func (result ExecutionResult) Empty() (empty bool) {
	return len(result.Diagnostics) == 0 &&
		result.Output == "" &&
		result.Command == CommandResult{}
}

func ViolationsFound() (err error) {
	return errors.New("violations found")
}
