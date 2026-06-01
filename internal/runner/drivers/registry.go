package drivers

import (
	"errors"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
)

var errViolationsFound = errors.New("violations found")

func Checkers() (registry runner.ExecutorRegistry) {
	return runner.ExecutorRegistry{
		contract.ExecutorToolchain:      runner.ToolchainExecutor,
		contract.ExecutorProject:        projectExecutor,
		contract.ExecutorFileCommand:    fileCommandExecutor,
		contract.ExecutorTargetCommand:  targetCommandExecutor,
		contract.ExecutorTargetCheck:    targetCheckExecutor,
		contract.ExecutorRepositoryScan: repositoryScanExecutor,
	}
}

func Fixers() (registry runner.ExecutorRegistry) {
	return runner.ExecutorRegistry{
		contract.ExecutorFileCommand:   fileCommandExecutor,
		contract.ExecutorTargetCommand: targetCommandExecutor,
	}
}
