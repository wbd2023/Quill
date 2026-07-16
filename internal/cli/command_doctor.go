package cli

import (
	"context"
	"flag"

	"ciphera/tools/internal/engine"
	"ciphera/tools/internal/report"
)

func runDoctor(tool Tool, options doctorOptions) (exitCode int) {
	doctor, err := engine.New(options.repoRoot)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	inspection, err := doctor.Inspect(context.Background())
	if err != nil {
		tool.writeError(err)
		return 1
	}

	result := report.ToolchainResult{Statuses: inspection.Statuses}
	if _, err = renderToolchainStatus(tool.stdout, options.format, result); err != nil {
		tool.writeError(err)
		return 1
	}

	if inspection.AllValid {
		return 0
	}

	return 1
}

func parseDoctorOptions(arguments []string) (options doctorOptions, err error) {
	return parseDoctorOptionsWithResolver(resolveRepoRoot, arguments)
}

func parseDoctorOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options doctorOptions, err error) {
	const summary = "check pinned style tools"
	var format string
	flagSet := newDoctorFlagSet(&options, &format)
	if err = parseArguments(flagSet, summary, arguments); err != nil {
		return options, err
	}

	options.format, err = parseFormat(format)
	if err != nil {
		return options, err
	}

	options.repoRoot, err = resolve(options.repoRoot)
	return options, err
}

func newDoctorFlagSet(options *doctorOptions, format *string) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("doctor")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	return flagSet
}

func doctorUsageText() (usage string) {
	const summary = "check pinned style tools"
	var options doctorOptions
	var format string
	return commandUsage("doctor", summary, newDoctorFlagSet(&options, &format))
}
