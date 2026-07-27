package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

const (
	helpCommand   = "help"
	usageExitCode = 2
)

type repositoryRootResolver func(string) (string, error)

// engineFactory builds the Engine used by a command. It defaults to engine.New and is overridable
// in white-box tests so a fake pack provider or cancellation behaviour can be injected without a
// real toolchain.
type engineFactory func(repositoryRoot string, options ...engine.Option) (*engine.Engine, error)

// Tool runs CLI commands using the configured output streams.
type Tool struct {
	stdout          io.Writer
	stderr          io.Writer
	resolveRepoRoot repositoryRootResolver
	version         string
	buildEngine     engineFactory
}

// Action is a parsed CLI command ready to run with the operation context.
type Action func(context.Context, Tool) int

// New constructs a CLI tool with the given output streams and build version.
func New(stdout io.Writer, stderr io.Writer, version string) (tool Tool) {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	return Tool{
		stdout:          stdout,
		stderr:          stderr,
		resolveRepoRoot: resolveRepoRoot,
		version:         version,
		buildEngine:     engine.New,
	}
}

// Run parses and executes one CLI command. The operation context is propagated to the engine and
// child tools; a cancelled context surfaces as a machine error envelope in JSON mode.
func (tool Tool) Run(ctx context.Context, arguments []string) (exitCode int) {
	if len(arguments) == 0 {
		tool.writeUsageError(rootUsageText(), nil)
		return usageExitCode
	}

	if isHelpRequest(arguments[0]) {
		return tool.runHelp(arguments[1:])
	}

	command, found := findCommand(arguments[0])
	if !found {
		tool.writeUsageError(rootUsageText(), fmt.Errorf("unknown command %q", arguments[0]))
		return usageExitCode
	}

	action, err := command.prepare(tool.resolveRepoRoot, arguments[1:])
	if err == nil {
		return action(ctx, tool)
	}

	var help flagHelpError
	if errors.As(err, &help) {
		_, _ = io.WriteString(tool.stdout, help.message)
		return 0
	}

	// A preparation failure is an invalid argument: malformed flags, an unexpected positional
	// argument, an unparseable format, or a repository-root resolution failure. In machine mode
	// emit a stable error envelope on stdout; otherwise write the human usage error to stderr.
	if command.machineMode(arguments[1:]) {
		tool.writeMachineErrorEnvelope(command.name, report.ErrorCodeInvalidArgument, err)
		return usageExitCode
	}

	tool.writeUsageError(command.usage(), err)
	return usageExitCode
}
