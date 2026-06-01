package pack

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/toolchain"
)

// Definition describes a modular checker collection.
type Definition struct {
	ID       string
	Name     string
	Tools    []toolchain.Capability
	Rules    []contract.RuleDefinition
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
func CloneRules(rules []contract.RuleDefinition) (clones []contract.RuleDefinition) {
	clones = make([]contract.RuleDefinition, 0, len(rules))
	for _, rule := range rules {
		clones = append(clones, cloneRule(rule))
	}

	return clones
}

func cloneRule(rule contract.RuleDefinition) (clone contract.RuleDefinition) {
	clone = rule
	clone.Check = cloneExecutionSpec(rule.Check)
	clone.Fix = cloneExecutionSpec(rule.Fix)
	return clone
}

func cloneExecutionSpec(spec contract.ExecutionSpec) (clone contract.ExecutionSpec) {
	clone = spec
	clone.Detail = cloneExecutionDetail(spec.Detail)
	return clone
}

func cloneExecutionDetail(detail contract.ExecutionDetail) (clone contract.ExecutionDetail) {
	switch execution := detail.(type) {
	case contract.ToolchainExecution:
		execution.ToolIDs = append([]string{}, execution.ToolIDs...)
		return execution

	case contract.FileCommandExecution:
		execution.Arguments = append([]string{}, execution.Arguments...)
		return execution

	case contract.TargetCommandExecution:
		execution.ToolIDs = append([]string{}, execution.ToolIDs...)
		execution.Targets = append([]string{}, execution.Targets...)
		return execution

	case contract.TargetCheckExecution:
		execution.ToolIDs = append([]string{}, execution.ToolIDs...)
		execution.Targets = append([]string{}, execution.Targets...)
		return execution

	default:
		return detail
	}
}
