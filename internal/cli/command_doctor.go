package cli

import (
	"context"
	"flag"

	"github.com/wbd2023/quill/internal/report"
)

func runDoctor(ctx context.Context, tool Tool, options doctorOptions) (exitCode int) {
	engineInstance, err := tool.buildEngine(options.repoRoot)
	if err != nil {
		return tool.reportCommandError(ctx, "doctor", options.format, err)
	}

	inspection, err := engineInstance.Inspect(ctx)
	if err != nil {
		return tool.reportCommandError(ctx, "doctor", options.format, err)
	}

	result := report.ToolchainResult{Statuses: inspection.Statuses}
	if _, err = renderToolchainStatus(tool.stdout, "doctor", options.format, result); err != nil {
		return tool.reportCommandError(ctx, "doctor", options.format, err)
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
