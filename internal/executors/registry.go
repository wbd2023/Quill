package executors

import (
	"errors"

	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

var errViolationsFound = errors.New("violations found")

func Checkers() (executors runner.ExecutorRegistry) {
	return runner.ExecutorRegistry{
		rulepack.ExecutorToolchain:      runner.ToolchainExecutor,
		rulepack.ExecutorControlPlane:   controlPlaneExecutor,
		rulepack.ExecutorFileCommand:    fileCommandExecutor,
		rulepack.ExecutorBackendCommand: backendCommandExecutor,
		rulepack.ExecutorBackendCheck:   backendCheckExecutor,
		rulepack.ExecutorRepositoryScan: repositoryScanExecutor,
	}
}

func Fixers() (executors runner.ExecutorRegistry) {
	return runner.ExecutorRegistry{
		rulepack.ExecutorFileCommand:    fileCommandExecutor,
		rulepack.ExecutorBackendCommand: backendCommandExecutor,
	}
}
