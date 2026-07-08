package toolchain

import "ciphera/tools/internal/style"

/* -------------------------------------------- Types ------------------------------------------- */

// Capability is a pinned external tool and how to inspect and install it.
type Capability struct {
	ID      string
	Name    string
	Command string
	Version VersionSpec
	Install InstallSpec
}

// VersionSpec selects how a tool's installed version is detected. Sealed: only the variants in
// this package satisfy it, dispatched by type-switch in detectVersion.
type VersionSpec interface {
	versionSpec()
}

// GoCommandVersion runs `go version` and parses the goX.Y.Z token.
type GoCommandVersion struct{}

// BuildInfoVersion reads embedded build info and checks the binary's module path.
type BuildInfoVersion struct {
	ModulePath string
}

// PrefixedLineVersion runs `<command> --version` and finds a "version:" prefixed line.
type PrefixedLineVersion struct{}

// FirstTokenVersion runs `<command> --version` and parses the first whitespace-delimited token.
type FirstTokenVersion struct{}

// InstallSpec selects how a missing tool is installed. Sealed: only the variants in this package
// satisfy it, dispatched by type-switch in installer.installTool.
type InstallSpec interface {
	installSpec()
}

// NoInstall means the tool is never installed by the engine (assumed present on the host).
type NoInstall struct{}

// GoBinaryInstall runs `go install <Source>@<version>`.
type GoBinaryInstall struct {
	Source string
}

// NodePackageInstall runs `npm install <Source>@<version>`.
type NodePackageInstall struct {
	Source string
}

// ArchiveInstall downloads, verifies, and extracts a release archive.
type ArchiveInstall struct {
	Spec ArchiveSpec
}

// ArchiveSpec describes how to download and extract a binary tool from a release archive. Carried
// on ArchiveInstall; both format strings are passed to fmt.Sprintf - URLFormat with args
// (version, platform) using indexed %[1]s for repeats, BinaryPathFormat with arg (version).
type ArchiveSpec struct {
	URLFormat        string
	BinaryPathFormat string
	Platforms        map[string]string
}

/* ------------------------------------------- Markers ------------------------------------------ */

func (GoCommandVersion) versionSpec()    {}
func (BuildInfoVersion) versionSpec()    {}
func (PrefixedLineVersion) versionSpec() {}
func (FirstTokenVersion) versionSpec()   {}

func (NoInstall) installSpec()          {}
func (GoBinaryInstall) installSpec()    {}
func (NodePackageInstall) installSpec() {}
func (ArchiveInstall) installSpec()     {}

/* ------------------------------------------- Helpers ------------------------------------------ */

func (capability Capability) Tool() (tool style.Tool) {
	return style.Tool{
		ID:   capability.ID,
		Name: capability.Name,
	}
}

// Policies converts capabilities to the style tools they represent.
func Policies(capabilities []Capability) (tools []style.Tool) {
	tools = make([]style.Tool, 0, len(capabilities))
	for _, capability := range capabilities {
		tools = append(tools, capability.Tool())
	}

	return tools
}

// CapabilitiesByID indexes tool capabilities by tool ID.
func CapabilitiesByID(
	capabilities []Capability,
) (indexed map[string]Capability) {
	indexed = make(map[string]Capability, len(capabilities))
	for _, capability := range capabilities {
		indexed[capability.ID] = capability
	}

	return indexed
}
