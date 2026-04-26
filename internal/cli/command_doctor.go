package cli

import (
	"flag"

	"ciphera/tools/internal/report"
)

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

func runDoctor(tool CLI, options doctorOptions) (exitCode int) {
	context, err := loadContext(options.repoRoot, "")
	if err != nil {
		tool.writeError(err)
		return 1
	}

	toolIDs := toolIDsFromTools(context.Effective.Tools)

	statuses, allValid := inspectToolchain(
		context.Effective.Tools,
		context.ToolCapabilities,
		toolIDs,
		context.ToolEnvironment,
	)
	result := report.ToolchainResult{Statuses: statuses}
	_, err = renderToolchainStatus(tool.stdout, options.format, result)
	if err != nil {
		tool.writeError(err)
		return 1
	}

	if allValid {
		return 0
	}

	return 1
}
