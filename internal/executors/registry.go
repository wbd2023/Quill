package executors

import (
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/runner"
)

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
