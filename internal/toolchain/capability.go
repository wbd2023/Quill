package toolchain

import "ciphera/tools/internal/contract"

type VersionKind string

type InstallKind string

type Capability struct {
	ID            string
	Name          string
	Command       string
	VersionKind   VersionKind
	ModulePath    string
	InstallKind   InstallKind
	InstallSource string
}

func (capability Capability) Tool() (tool contract.Tool) {
	return contract.Tool{
		ID:   capability.ID,
		Name: capability.Name,
	}
}

func Policies(capabilities []Capability) (tools []contract.Tool) {
	tools = make([]contract.Tool, 0, len(capabilities))
	for _, capability := range capabilities {
		tools = append(tools, capability.Tool())
	}

	return tools
}

func CapabilitiesByID(
	capabilities []Capability,
) (indexed map[string]Capability) {
	indexed = make(map[string]Capability, len(capabilities))
	for _, capability := range capabilities {
		indexed[capability.ID] = capability
	}

	return indexed
}
