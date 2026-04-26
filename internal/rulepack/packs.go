package rulepack

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/toolchain"
)

type Pack struct {
	ID    string
	Name  string
	Tools []toolchain.Capability
	Rules []RuleDefinition
}

type RuleDefinition = contract.RuleDefinition

type ExecutionSpec = contract.ExecutionSpec

type Registry struct {
	packs        []Pack
	capabilities []toolchain.Capability
	rules        []RuleDefinition
}

func (registry Registry) Packs() (packs []Pack) {
	return append([]Pack{}, registry.packs...)
}

func (registry Registry) ToolCapabilities() (capabilities []toolchain.Capability) {
	return append([]toolchain.Capability{}, registry.capabilities...)
}

func (registry Registry) Tools() (tools []contract.Tool) {
	return toolchain.Policies(registry.capabilities)
}

func (registry Registry) Rules() (rules []RuleDefinition) {
	return append([]RuleDefinition{}, registry.rules...)
}

func (registry Registry) Definitions() (definitions contract.Definitions) {
	return contract.Definitions{
		Tools: registry.Tools(),
		Rules: registry.Rules(),
	}
}

func (registry Registry) ToolByID(id string) (capability toolchain.Capability, found bool) {
	for _, capability := range registry.capabilities {
		if capability.ID == id {
			return capability, true
		}
	}

	return toolchain.Capability{}, false
}
