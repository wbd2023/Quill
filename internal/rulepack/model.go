package rulepack

import "ciphera/tools/internal/contract"

/* ------------------------------------------ Constants ----------------------------------------- */

const (
	PackControl    = "control"
	PackText       = "text"
	PackMarkdown   = "markdown"
	PackShell      = "shell"
	PackGo         = "go"
	PackRepository = "repository"
)

const (
	ConfigRefArchitecture = "architecture"
	ConfigRefControlPlane = "control_plane"
	ConfigRefNaming       = "naming"
	ConfigRefRepository   = "repository"
)

const (
	LanguageGo = "go"
)

/* -------------------------------------------- Types ------------------------------------------- */

type Pack struct {
	ID    string
	Name  string
	Tools []contract.Tool
	Rules []RuleDefinition
}

type RuleDefinition = contract.RuleDefinition

type ExecutionSpec = contract.ExecutionSpec

type Registry struct {
	packs []Pack
	tools []contract.Tool
	rules []RuleDefinition
}

/* ------------------------------------------ Accessors ----------------------------------------- */

func (registry Registry) Packs() (packs []Pack) {
	return append([]Pack{}, registry.packs...)
}

func (registry Registry) Tools() (tools []contract.Tool) {
	return append([]contract.Tool{}, registry.tools...)
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

func (registry Registry) ToolByID(id string) (tool contract.Tool, found bool) {
	for _, current := range registry.tools {
		if current.ID == id {
			return current, true
		}
	}

	return contract.Tool{}, false
}
