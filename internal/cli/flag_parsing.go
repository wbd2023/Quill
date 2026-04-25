package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/report"
)

/* ------------------------------------------ Flag Sets ----------------------------------------- */

func newFlagSet(name string) (flagSet *flag.FlagSet) {
	flagSet = flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flagSet.Usage = func() {}
	return flagSet
}

/* -------------------------------------- Argument Parsing -------------------------------------- */

func parseArguments(flagSet *flag.FlagSet, summary string, arguments []string) (err error) {
	if err = flagSet.Parse(arguments); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return flagHelpError{
				message: commandUsage(flagSet.Name(), summary, flagSet),
			}
		}

		return err
	}

	return ensureNoPositionalArguments(flagSet.Args())
}

/* ---------------------------------------- Value Parsing --------------------------------------- */

func parseScope(value string) (scope contract.Scope, err error) {
	switch contract.Scope(value) {
	case contract.ScopeApp, contract.ScopeTools, contract.ScopeAll:
		return contract.Scope(value), nil
	default:
		return contract.ScopeUnknown, fmt.Errorf(
			"invalid scope %q: must be app, tools, or all",
			value,
		)
	}
}

func parseProfile(value string) (profile contract.CheckProfile, err error) {
	switch contract.CheckProfile(value) {
	case contract.CheckProfileRequired, contract.CheckProfileAll:
		return contract.CheckProfile(value), nil
	default:
		return "", fmt.Errorf("invalid profile %q: must be required or all", value)
	}
}

func parseFormat(value string) (format report.OutputFormat, err error) {
	switch report.OutputFormat(value) {
	case report.FormatText, report.FormatJSON:
		return report.OutputFormat(value), nil
	default:
		return "", fmt.Errorf("invalid format %q: must be text or json", value)
	}
}

/* ------------------------------------ Positional Arguments ------------------------------------ */

func ensureNoPositionalArguments(arguments []string) (err error) {
	if len(arguments) == 0 {
		return nil
	}

	return fmt.Errorf(
		"unexpected positional arguments: %s",
		strings.Join(arguments, ", "),
	)
}
