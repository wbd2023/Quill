package projectpolicy

// config constants.
const (
	// CommandsRunnerMake means quality commands are declared through Make targets.
	CommandsRunnerMake CommandsRunner = "make"
)

// Config defines Project Pack Policy.
type Config struct {
	Commands CommandsConfig
}

// CommandsRunner identifies how a repository exposes its quality commands.
type CommandsRunner string

// CommandsConfig describes the repository quality commands expected by project checks.
type CommandsConfig struct {
	Runner CommandsRunner
	Make   MakeConfig
}

// MakeConfig describes Make-backed quality commands.
type MakeConfig struct {
	Path              string
	RequiredVariables []MakefileVariable
	RequiredTargets   []MakefileTarget
}

// MakefileVariable describes a required Makefile variable.
type MakefileVariable struct {
	Name  string
	Value string
}

// MakefileTarget describes a required Makefile target.
type MakefileTarget struct {
	Name       string
	RecipeLine string
}
