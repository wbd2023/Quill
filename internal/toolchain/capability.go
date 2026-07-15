package toolchain

// Capability describes an external tool's command and the strategies for detecting its version and
// installing it.
type Capability struct {
	ID      string
	Name    string
	Command string

	Version VersionMethod
	Install InstallMethod
}

// VersionMethod selects how a tool's installed version is detected.
type VersionMethod interface {
	versionMethod()
}

// GoVersion runs `go version` and parses the goX.Y.Z token.
type GoVersion struct{}

// ModuleVersion reads embedded build info; ModulePath, if set, must match the binary's main module.
type ModuleVersion struct {
	ModulePath string
}

// PrefixedLineVersion runs `<command> --version` and finds a "version:" prefixed line.
type PrefixedLineVersion struct{}

// FirstTokenVersion runs `<command> --version` and parses the first whitespace-delimited token.
type FirstTokenVersion struct{}

// InstallMethod selects how a missing tool is installed.
type InstallMethod interface {
	installMethod()
}

// NoInstall means the tool is never installed by the engine (assumed present on the host).
type NoInstall struct{}

// GoInstall runs `go install <Source>@<version>`.
type GoInstall struct {
	Source string
}

// NpmInstall runs `npm install <Source>@<version>`.
type NpmInstall struct {
	Source string
}

// GitHubInstall installs a tool from a GitHub release archive.
type GitHubInstall struct {
	Owner      string
	Repository string

	// Tag is a fmt.Sprintf format string that produces the release tag, taking the version.
	Tag string

	// Asset is a fmt.Sprintf format string that produces the archive filename, taking the tag and
	// platform.
	Asset string

	// Path is a fmt.Sprintf format string that produces the executable path inside the archive,
	// taking the tag.
	Path string

	// Platforms maps "os/arch" to the per-platform token substituted into Asset.
	Platforms map[string]string
}

func (GoVersion) versionMethod()           {}
func (ModuleVersion) versionMethod()       {}
func (PrefixedLineVersion) versionMethod() {}
func (FirstTokenVersion) versionMethod()   {}

func (NoInstall) installMethod()     {}
func (GoInstall) installMethod()     {}
func (NpmInstall) installMethod()    {}
func (GitHubInstall) installMethod() {}
