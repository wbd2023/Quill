package cli

import "context"

// Command is one supported CLI command and its preparation contract.
type Command struct {
	name        string
	summary     string
	usage       func() string
	prepare     func(repositoryRootResolver, []string) (Action, error)
	machineMode func([]string) bool
}

/* -------------------------------------- Command Registry -------------------------------------- */

var commands = []Command{
	{
		name:        "check",
		summary:     "run STYLE.md checks",
		usage:       checkUsageText,
		prepare:     prepareAction(parseCheckOptionsWithResolver, runCheck),
		machineMode: checkMachineMode,
	},
	{
		name:        "fix",
		summary:     "run safe style auto-fixes",
		usage:       fixUsageText,
		prepare:     prepareAction(parseFixOptionsWithResolver, runFix),
		machineMode: fixMachineMode,
	},
	{
		name:        "doctor",
		summary:     "check pinned style tools",
		usage:       doctorUsageText,
		prepare:     prepareAction(parseDoctorOptionsWithResolver, runDoctor),
		machineMode: doctorMachineMode,
	},
	{
		name:        "coverage",
		summary:     "show STYLE.md automation coverage",
		usage:       coverageUsageText,
		prepare:     prepareAction(parseCoverageOptionsWithResolver, runCoverage),
		machineMode: coverageMachineMode,
	},
	{
		name:        "install",
		summary:     "install pinned style tools",
		usage:       installUsageText,
		prepare:     prepareAction(parseInstallOptionsWithResolver, runInstall),
		machineMode: installMachineMode,
	},
	{
		name:        "lock",
		summary:     "resolve archive-tool hashes to quill.lock",
		usage:       lockUsageText,
		prepare:     prepareAction(parseLockOptionsWithResolver, runLock),
		machineMode: lockMachineMode,
	},
	{
		name:        "version",
		summary:     "print the Quill version",
		usage:       versionUsageText,
		prepare:     prepareAction(parseVersionOptions, runVersion),
		machineMode: versionMachineMode,
	},
	{
		name:        "init",
		summary:     "create a minimal STYLE.md and quill.toml",
		usage:       initUsageText,
		prepare:     prepareAction(parseInitOptionsWithResolver, runInit),
		machineMode: initMachineMode,
	},
	{
		name:        "list",
		summary:     "list packs, rules, tools, or scopes",
		usage:       listUsageText,
		prepare:     prepareAction(parseListOptionsWithResolver, runList),
		machineMode: listMachineMode,
	},
	{
		name:        "explain",
		summary:     "explain an active rule",
		usage:       explainUsageText,
		prepare:     prepareAction(parseExplainOptionsWithResolver, runExplain),
		machineMode: explainMachineMode,
	},
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func findCommand(name string) (matched Command, found bool) {
	for _, command := range commands {
		if command.name == name {
			return command, true
		}
	}

	return Command{}, false
}

func prepareAction[options any](
	parse func(repositoryRootResolver, []string) (options, error),
	run func(context.Context, Tool, options) int,
) (prepare func(repositoryRootResolver, []string) (Action, error)) {
	return func(resolve repositoryRootResolver, arguments []string) (bound Action, err error) {
		options, err := parse(resolve, arguments)
		if err != nil {
			return nil, err
		}

		return func(ctx context.Context, tool Tool) int {
			return run(ctx, tool, options)
		}, nil
	}
}
