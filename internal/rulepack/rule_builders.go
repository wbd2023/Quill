package rulepack

import "ciphera/tools/internal/contract"

/* ---------------------------------------- Rule Builders --------------------------------------- */

func toolchainRule(
	id string,
	name string,
	toolIDs ...string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:    id,
		Name:  name,
		Group: contract.RuleGroupControlPlane,
		Spec: contract.ExecutionSpec{
			Executor: contract.ExecutorToolchain,
			ToolIDs:  append([]string{}, toolIDs...),
		},
	}
}

func controlPlaneRule(
	id string,
	name string,
	check string,
) (rule RuleDefinition) {
	return controlPlaneRuleWithConfig(id, name, check, ConfigRefControlPlane)
}

func controlPlaneRuleWithConfig(
	id string,
	name string,
	check string,
	configRef string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:                 id,
		Name:               name,
		Group:              contract.RuleGroupControlPlane,
		RequiredConfigRefs: []string{configRef},
		Spec: contract.ExecutionSpec{
			Executor: contract.ExecutorControlPlane,
			Check:    check,
		},
	}
}

func fileCommandRule(
	id string,
	name string,
	toolID string,
	fileSet string,
	arguments []string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:    id,
		Name:  name,
		Group: contract.RuleGroupExternal,
		Spec: contract.ExecutionSpec{
			Executor:  contract.ExecutorFileCommand,
			ToolID:    toolID,
			FileSet:   fileSet,
			Arguments: append([]string{}, arguments...),
		},
	}
}

func fileCommandRuleWithConfig(
	id string,
	name string,
	toolID string,
	fileSet string,
	arguments []string,
	configArgument string,
	configFile string,
) (rule RuleDefinition) {
	rule = fileCommandRule(id, name, toolID, fileSet, arguments)
	rule.Spec.ConfigArgument = configArgument
	rule.Spec.ConfigFile = configFile
	return rule
}

func golangciRule(
	id string,
	name string,
	backend string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:    id,
		Name:  name,
		Group: contract.RuleGroupLanguage,
		Spec: contract.ExecutionSpec{
			Executor: contract.ExecutorGolangci,
			ToolIDs: []string{
				contract.ToolGo,
				contract.ToolGoimports,
				contract.ToolGolangciLint,
			},
			Backend:  backend,
			Language: LanguageGo,
		},
		FixSpec: contract.ExecutionSpec{
			Executor: contract.ExecutorGoFormat,
			ToolIDs: []string{
				contract.ToolGo,
				contract.ToolGoimports,
			},
			Backend:  backend,
			Language: LanguageGo,
		},
	}
}

func goStyleRule(
	id string,
	name string,
	backend string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:                  id,
		Name:                name,
		Group:               contract.RuleGroupLanguage,
		RequiredPathClasses: requiredGoStylePathClasses(),
		Spec: contract.ExecutionSpec{
			Executor: contract.ExecutorGoStyle,
			ToolIDs:  []string{contract.ToolGo},
			Backend:  backend,
			Language: LanguageGo,
		},
	}
}

func repoScanRule(
	id string,
	name string,
	scanner string,
) (rule RuleDefinition) {
	return repoScanRuleWithConfig(id, name, scanner)
}

func repoScanRuleWithConfig(
	id string,
	name string,
	scanner string,
	configRefs ...string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:                 id,
		Name:               name,
		Group:              contract.RuleGroupRepository,
		RequiredConfigRefs: append([]string{}, configRefs...),
		Spec: contract.ExecutionSpec{
			Executor: contract.ExecutorRepositoryScan,
			Scanner:  scanner,
		},
	}
}

func lineLengthRule() (rule RuleDefinition) {
	rule = repoScanRule(
		"repo/line-length",
		"Line length",
		RepositoryScannerLineLength,
	)
	rule.Spec.FileSet = "line_length"
	return rule
}
