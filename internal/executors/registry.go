package executors

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runner"
)

func Checkers() (executors runner.ExecutorRegistry) {
	return runner.ExecutorRegistry{
		contract.ExecutorToolchain:      runner.ToolchainExecutor,
		contract.ExecutorControlPlane:   controlPlaneExecutor,
		contract.ExecutorFileCommand:    runner.FileCommandExecutor,
		contract.ExecutorGolangci:       golangciExecutor,
		contract.ExecutorGoStyle:        goStyleExecutor,
		contract.ExecutorRepositoryScan: repositoryScanExecutor,
	}
}

func Fixers() (executors runner.ExecutorRegistry) {
	return runner.ExecutorRegistry{
		contract.ExecutorFileCommand: runner.FileCommandExecutor,
		contract.ExecutorGoFormat:    goFormatExecutor,
	}
}
