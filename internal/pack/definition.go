package pack

import (
	"slices"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

// Definition describes a modular checker collection.
type Definition struct {
	ID       string
	Name     string
	Tools    []toolchain.Capability
	Rules    []style.RuleDefinition
	FileSets policy.FileSets
	Config   Config
}

// Config describes the profile config accepted by a pack.
type Config struct {
	Required bool
	Validate func(policy.PackConfig) error
}

// CloneDefinitions returns deep copies of the supplied pack definitions.
func CloneDefinitions(definitions []Definition) (clones []Definition) {
	clones = make([]Definition, 0, len(definitions))
	for _, definition := range definitions {
		clones = append(clones, CloneDefinition(definition))
	}

	return clones
}

// CloneDefinition returns a deep copy of definition.
func CloneDefinition(definition Definition) (clone Definition) {
	clone = definition
	clone.Tools = append([]toolchain.Capability{}, definition.Tools...)
	clone.Rules = CloneRules(definition.Rules)
	clone.FileSets = definition.FileSets.Clone()
	return clone
}

// CloneRules returns deep copies of the supplied rule definitions.
func CloneRules(rules []style.RuleDefinition) (clones []style.RuleDefinition) {
	clones = make([]style.RuleDefinition, 0, len(rules))
	for _, rule := range rules {
		clones = append(clones, cloneRule(rule))
	}

	return clones
}

func cloneRule(rule style.RuleDefinition) (clone style.RuleDefinition) {
	clone = rule
	clone.Check = cloneTemplate(rule.Check)
	clone.Fix = cloneTemplate(rule.Fix)
	return clone
}
func cloneTemplate(template style.Template) (clone style.Template) {
	switch detail := template.(type) {
	case style.ToolchainExecution:
		detail.ToolIDs = slices.Clone(detail.ToolIDs)
		return detail

	case style.FileCommandExecution:
		detail.Arguments = slices.Clone(detail.Arguments)
		return detail

	case style.TargetCommandTemplate:
		detail.ToolIDs = slices.Clone(detail.ToolIDs)
		return detail

	case style.TargetCheckTemplate:
		detail.ToolIDs = slices.Clone(detail.ToolIDs)
		return detail

	default:
		return template
	}
}
