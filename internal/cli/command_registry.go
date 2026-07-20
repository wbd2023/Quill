package cli

// Command is command.
type Command struct {
	name    string
	summary string
	usage   func() string
	prepare func(repositoryRootResolver, []string) (Action, error)
}

var commands = []Command{
	{
		name:    "check",
		summary: "run STYLE.md checks",
		usage:   checkUsageText,
		prepare: prepareAction(parseCheckOptionsWithResolver, runCheck),
	},
	{
		name:    "fix",
		summary: "run safe style auto-fixes",
		usage:   fixUsageText,
		prepare: prepareAction(parseFixOptionsWithResolver, runFix),
	},
	{
		name:    "doctor",
		summary: "check pinned style tools",
		usage:   doctorUsageText,
		prepare: prepareAction(parseDoctorOptionsWithResolver, runDoctor),
	},
	{
		name:    "coverage",
		summary: "show STYLE.md automation coverage",
		usage:   coverageUsageText,
		prepare: prepareAction(parseCoverageOptionsWithResolver, runCoverage),
	},
	{
		name:    "install",
		summary: "install pinned style tools",
		usage:   installUsageText,
		prepare: prepareAction(parseInstallOptionsWithResolver, runInstall),
	},
	{
		name:    "lock",
		summary: "resolve archive-tool hashes to quill.lock",
		usage:   lockUsageText,
		prepare: prepareAction(parseLockOptionsWithResolver, runLock),
	},
	{
		name:    "version",
		summary: "print the Quill version",
		usage:   versionUsageText,
		prepare: prepareAction(parseVersionOptions, runVersion),
	},
}

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
	run func(Tool, options) int,
) (prepare func(repositoryRootResolver, []string) (Action, error)) {
	return func(resolve repositoryRootResolver, arguments []string) (bound Action, err error) {
		options, err := parse(resolve, arguments)
		if err != nil {
			return nil, err
		}

		return func(tool Tool) int {
			return run(tool, options)
		}, nil
	}
}
