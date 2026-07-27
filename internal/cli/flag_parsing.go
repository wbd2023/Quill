package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/wbd2023/quill/internal/report"
	"github.com/wbd2023/quill/internal/style"
)

/* -------------------------------------- Flag Construction ------------------------------------- */

func newFlagSet(name string) (flagSet *flag.FlagSet) {
	flagSet = flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.SetOutput(io.Discard)
	flagSet.Usage = func() {}
	return flagSet
}

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

func parseScope(value string) (scope style.Scope, err error) {
	if strings.TrimSpace(value) == "" {
		return "", nil
	}

	return style.Scope(value), nil
}

func parseCheckMode(value string) (mode style.CheckMode, err error) {
	switch style.CheckMode(value) {
	case style.CheckModeRequired, style.CheckModeAll:
		return style.CheckMode(value), nil
	default:
		return "", fmt.Errorf("invalid mode %q: must be required or all", value)
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

/* ------------------------------------- Argument Validation ------------------------------------ */

func ensureNoPositionalArguments(arguments []string) (err error) {
	if len(arguments) == 0 {
		return nil
	}

	return fmt.Errorf(
		"unexpected positional arguments: %s",
		strings.Join(arguments, ", "),
	)
}

/* --------------------------------------- Machine Output --------------------------------------- */

// machineModeRequested reports whether the command's flag parser reaches `--format json` before
// it stops or fails. It parses raw arguments so a preparation failure can still render a machine
// error envelope using the format that the command itself established.
func machineModeRequested(
	arguments []string,
	flagSet *flag.FlagSet,
	format *string,
) (requested bool) {
	_ = flagSet.Parse(arguments)
	return report.OutputFormat(*format) == report.FormatJSON
}

func checkMachineMode(arguments []string) (requested bool) {
	var options checkOptions
	var scope string
	var mode string
	var format string
	return machineModeRequested(
		arguments, newCheckFlagSet(&options, &scope, &mode, &format), &format,
	)
}

func fixMachineMode(arguments []string) (requested bool) {
	var options fixOptions
	var scope string
	var format string
	return machineModeRequested(arguments, newFixFlagSet(&options, &scope, &format), &format)
}

func doctorMachineMode(arguments []string) (requested bool) {
	var options doctorOptions
	var format string
	return machineModeRequested(arguments, newDoctorFlagSet(&options, &format), &format)
}

func coverageMachineMode(arguments []string) (requested bool) {
	var options coverageOptions
	var format string
	return machineModeRequested(arguments, newCoverageFlagSet(&options, &format), &format)
}

func installMachineMode(arguments []string) (requested bool) {
	var options installOptions
	var format string
	return machineModeRequested(arguments, newInstallFlagSet(&options, &format), &format)
}

func lockMachineMode(arguments []string) (requested bool) {
	var options lockOptions
	var format string
	return machineModeRequested(arguments, newLockFlagSet(&options, &format), &format)
}

func versionMachineMode(arguments []string) (requested bool) {
	return machineModeRequested(arguments, newFlagSet("version"), new(string))
}
