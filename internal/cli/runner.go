package cli

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/alecthomas/kong"

	"github.com/wbd2023/quill/internal/report"
)

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	helpCommand                    = "help"
	documentedHelpMaximumArguments = 2
	usageExitCode                  = 2
)

/* -------------------------------------------- Types ------------------------------------------- */

// Runner is the invocation-local CLI adapter. It owns command parsing, output streams, build
// metadata, and command lifecycle.
type Runner struct {
	stdout  io.Writer
	stderr  io.Writer
	version string
}

/* ----------------------------------------- Constructor ---------------------------------------- */

// New constructs a CLI runner with the given output streams and build version.
func New(stdout io.Writer, stderr io.Writer, version string) (runner Runner) {
	if stdout == nil {
		stdout = io.Discard
	}
	if stderr == nil {
		stderr = io.Discard
	}

	return Runner{
		stdout:  stdout,
		stderr:  stderr,
		version: version,
	}
}

/* --------------------------------------------- Run -------------------------------------------- */

// Run parses and executes one Quill CLI invocation.
//
// Kong owns grammar, parsing, defaults, static enum and positional validation, and command
// selection. Quill owns output policy, deterministic help, JSON envelopes, stable error codes,
// repository and target preparation, and 0/1/2 exit statuses.
//
// Help is handled before preparation and execution. A help flag hook returns a private sentinel
// wrapped by Kong in a ParseError; Run detects it and writes deterministic help to stdout with exit
// status 0.
//
// Parse failures are usage errors. Text mode writes usage to stderr; JSON mode writes an
// invalid_argument envelope to stdout. Parse errors use the output format from the original parse
// context rather than reparsing the arguments.
func (r Runner) Run(ctx context.Context, arguments []string) (exitCode int) {
	parser, model, err := newParser()
	if err != nil {
		r.writeError(err)
		return 1
	}

	if target, argument, malformed := malformedHelpAlias(arguments); malformed {
		if node := commandNode(parser.Model, target); node != nil {
			r.writeUsageError(commandUsage(node), fmt.Errorf("unexpected argument %q", argument))
			return usageExitCode
		}
	}

	arguments = translateHelpCommand(arguments)

	kctx, parseErr := parser.Parse(arguments)
	if parseErr != nil {
		if errors.Is(parseErr, errHelpRequested) {
			r.writeHelp(parser, selectedNode(parseErrorContext(parseErr)))
			return 0
		}
		return r.reportParseError(parser, parseErr)
	}

	selected := kctx.Selected()
	if selected == nil {
		r.writeUsageError(rootUsage(parser.Model), nil)
		return usageExitCode
	}

	cmd := model.lookup(selected)
	if cmd == nil {
		r.writeUsageError(rootUsage(parser.Model), nil)
		return usageExitCode
	}

	return cmd.run(ctx, r)
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// translateHelpCommand rewrites the documented `help` / `help <command>` syntax into the
// `<command> --help` flag form so Kong's single parser handles every help surface. `help` alone
// becomes `--help` (root help); `help check` becomes `check --help` (command help).
func translateHelpCommand(arguments []string) (translated []string) {
	if len(arguments) == 0 || arguments[0] != helpCommand {
		return arguments
	}

	switch len(arguments) {
	case 1:
		return []string{"--help"}
	case documentedHelpMaximumArguments:
		return []string{arguments[1], "--help"}
	default:
		return arguments
	}
}

// malformedHelpAlias reports the first unexpected argument when the documented help alias names
// a command and supplies additional tokens. It lets Runner retain command-specific usage without
// teaching Kong a second command grammar.
func malformedHelpAlias(arguments []string) (target string, argument string, malformed bool) {
	if len(arguments) <= documentedHelpMaximumArguments || arguments[0] != helpCommand {
		return "", "", false
	}

	return arguments[1], arguments[2], true
}

func commandNode(model *kong.Application, name string) (command *kong.Node) {
	for _, node := range model.Children {
		if node.Name == name {
			return node
		}
	}

	return nil
}

// reportParseError classifies a parse failure by the format established on the command-line. The
// format and selected command are extracted from the ParseError context rather than by reparsing
// the raw arguments.
func (r Runner) reportParseError(parser *kong.Kong, err error) (exitCode int) {
	context := parseErrorContext(err)
	node := selectedNode(context)

	if outputFormatFromContext(context) == report.FormatJSON {
		command := ""
		if node != nil {
			command = node.Name
		}
		r.writeMachineErrorEnvelope(
			command,
			report.ErrorCodeInvalidArgument,
			err,
		)
		return usageExitCode
	}

	usage := rootUsage(parser.Model)
	if node != nil {
		usage = commandUsage(node)
	}
	r.writeUsageError(usage, err)
	return usageExitCode
}

func (r Runner) writeHelp(parser *kong.Kong, node *kong.Node) {
	if node == nil {
		_, _ = io.WriteString(r.stdout, rootUsage(parser.Model))
		return
	}
	_, _ = io.WriteString(r.stdout, commandUsage(node))
}

// parseErrorContext returns the Kong parse context carried by a ParseError, or nil when the error
// is not a ParseError (Kong wraps every Parse failure in a ParseError, so this is defensive).
func parseErrorContext(err error) (context *kong.Context) {
	var parseErr *kong.ParseError
	if errors.As(err, &parseErr) {
		return parseErr.Context
	}
	return nil
}

func selectedNode(context *kong.Context) (node *kong.Node) {
	if context == nil {
		return nil
	}
	return context.Selected()
}

// outputFormatFromContext reports the output format established on the command-line at the point
// of failure. Flags accumulate their parsed values during tracing, so even a failed parse exposes
// a `--format json` that preceded the error. An absent format flag resolves to the text default.
func outputFormatFromContext(context *kong.Context) (format report.OutputFormat) {
	if context == nil {
		return report.FormatText
	}
	for _, flag := range context.Flags() {
		if flag.Name != "format" {
			continue
		}
		value := context.FlagValue(flag)
		if value == nil {
			return report.FormatText
		}
		return report.OutputFormat(fmt.Sprintf("%v", value))
	}
	return report.FormatText
}
