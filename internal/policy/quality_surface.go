package policy

const (
	// QualitySurfaceDriverMake means quality commands are declared through Make targets.
	QualitySurfaceDriverMake QualitySurfaceDriver = "make"
)

// QualitySurfaceDriver identifies how a repository exposes its quality commands.
type QualitySurfaceDriver string

// QualitySurfaceConfig defines how a repository exposes style quality commands.
type QualitySurfaceConfig struct {
	Driver QualitySurfaceDriver
	Make   MakeConfig
}

// MakeConfig defines the Make-specific quality command contract.
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
