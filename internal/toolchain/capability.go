package toolchain

import "ciphera/tools/internal/style"

// VersionKind is version kind.
type VersionKind string

// InstallKind is install kind.
type InstallKind string

// Capability is capability.
type Capability struct {
	ID            string
	Name          string
	Command       string
	VersionKind   VersionKind
	ModulePath    string
	InstallKind   InstallKind
	InstallSource string
}

func (capability Capability) Tool() (tool style.Tool) {
	return style.Tool{
		ID:   capability.ID,
		Name: capability.Name,
	}
}

// Policies returns the requested value.
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
