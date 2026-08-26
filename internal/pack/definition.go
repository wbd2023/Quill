package pack

import (
	"slices"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
)

// Definition describes a modular checker collection. Tools reference the canonical Tool
// capabilities owned by the catalog by global ID rather than carrying copies; the catalog
// resolves each reference and rejects unknown or duplicate declarations.
type Definition struct {
	ID       string
	Name     string
	ToolIDs  []string
	Rules    []style.RuleDefinition
	FileSets profile.FileSets
	Policy   Policy
}

// Policy describes the profile policy accepted by a Pack.
type Policy struct {
	Required bool
	Validate func(profile.PackPolicy) error
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
	clone.ToolIDs = slices.Clone(definition.ToolIDs)
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

// cloneTemplate returns a deep copy of an execution Template, cloning every slice the Template
// owns so a Pack's declarations stay isolated when cloned for the registry.
func cloneTemplate(template style.Template) (clone style.Template) {
	switch detail := template.(type) {
	case style.ToolchainCheck:
		detail.ToolIDs = slices.Clone(detail.ToolIDs)
		return detail

	case style.FileCommand:
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
