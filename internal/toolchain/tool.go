package toolchain

// Tool is a capability joined with its pinned version and execution limits, ready for inspection,
// installation, or execution.
type Tool struct {
	ID   string
	Name string

	PinnedVersion    string
	TimeoutSeconds   int
	OutputLimitBytes int64

	Command string
	Version VersionMethod
	Install InstallMethod
}
