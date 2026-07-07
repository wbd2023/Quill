package toolchain

import "ciphera/tools/internal/style"

// Known version-detection strategies, dispatched by detectVersion.
const (
	VersionKindGoCommand  VersionKind = "go_command"
	VersionKindBuildInfo  VersionKind = "build_info"
	VersionKindShellcheck VersionKind = "shellcheck"
	VersionKindNodeCLI    VersionKind = "node_cli"
)

// Known install strategies, dispatched by installer.installTool. InstallKindNone is a no-op
// success, distinct from an unset InstallKind which is rejected as unsupported.
const (
	InstallKindNone              InstallKind = "none"
	InstallKindGoBinary          InstallKind = "go_binary"
	InstallKindNodePackage       InstallKind = "node_package"
	InstallKindShellcheckArchive InstallKind = "shellcheck_archive"
)

// Known archive compression formats.
const (
	ArchiveFormatXz ArchiveFormat = "xz"
)

// VersionKind selects how a tool's installed version is detected.
type VersionKind string

// InstallKind selects how a missing tool is installed.
type InstallKind string

// ArchiveFormat selects the compression format of a release archive.
type ArchiveFormat string

// ArchiveSpec describes how to download and extract a binary tool from a release archive.
// Carried on Capability; nil for non-archive tools.
type ArchiveSpec struct {
	URL        func(version string, platform string) string
	Format     ArchiveFormat
	BinaryPath func(version string) string
	Platforms  map[string]string
}

// Capability is a pinned external tool and how to inspect and install it.
type Capability struct {
	ID            string
	Name          string
	Command       string
	VersionKind   VersionKind
	ModulePath    string
	InstallKind   InstallKind
	InstallSource string
	Archive       *ArchiveSpec
}

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
