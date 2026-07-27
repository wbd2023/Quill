package cli

import (
	"context"
	"flag"
	"fmt"

	"github.com/wbd2023/quill/internal/engine"
	"github.com/wbd2023/quill/internal/report"
)

func runLock(ctx context.Context, tool Tool, options lockOptions) (exitCode int) {
	progressWriter := tool.stdout
	if options.format == report.FormatJSON {
		// Machine mode reserves stdout for the single envelope; route lock progress to stderr.
		progressWriter = tool.stderr
	}

	engineInstance, err := tool.buildEngine(
		options.repoRoot, engine.WithProgressWriter(progressWriter),
	)
	if err != nil {
		return tool.reportCommandError(ctx, "lock", options.format, err)
	}

	result, err := engineInstance.Lock(ctx)
	if err != nil {
		return tool.reportCommandError(ctx, "lock", options.format, err)
	}

	if options.format == report.FormatJSON {
		return writeLockEnvelope(tool, "lock", result)
	}

	if _, err = fmt.Fprintf(
		tool.stdout, "Wrote %s (%d tools)\n", result.Path, result.ArchiveCount,
	); err != nil {
		tool.writeError(err)
		return 1
	}

	return 0
}

func writeLockEnvelope(tool Tool, command string, result engine.LockResult) (exitCode int) {
	if err := report.WriteLock(tool.stdout, command, report.LockResult{
		Path:         result.Path,
		ArchiveCount: result.ArchiveCount,
	}); err != nil {
		tool.writeError(err)
		return 1
	}

	return 0
}

func parseLockOptionsWithResolver(
	resolve repositoryRootResolver,
	arguments []string,
) (options lockOptions, err error) {
	const summary = "resolve and write archive-tool hashes to quill.lock"
	var format string
	flagSet := newLockFlagSet(&options, &format)
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

func newLockFlagSet(options *lockOptions, format *string) (flagSet *flag.FlagSet) {
	flagSet = newFlagSet("lock")
	flagSet.StringVar(
		&options.repoRoot,
		"repo-root",
		"",
		"repository root (auto-detected when omitted)",
	)
	flagSet.StringVar(format, "format", string(report.FormatText), "format: text|json")
	return flagSet
}

func lockUsageText() (usage string) {
	const summary = "resolve and write archive-tool hashes to quill.lock"
	var options lockOptions
	var format string
	return commandUsage("lock", summary, newLockFlagSet(&options, &format))
}
