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
		Group: RuleGroupControlPlane,
		Spec: contract.ExecutionSpec{
			Kind: ExecutorToolchain,
			Detail: contract.ToolchainExecution{
				ToolIDs: append([]string{}, toolIDs...),
			},
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
		Group:              RuleGroupControlPlane,
		RequiredConfigRefs: []string{configRef},
		Spec: contract.ExecutionSpec{
			Kind: ExecutorControlPlane,
			Detail: contract.ControlPlaneExecution{
				Check: check,
			},
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
		Group: RuleGroupExternal,
		Spec: contract.ExecutionSpec{
			Kind: ExecutorFileCommand,
			Detail: contract.FileCommandExecution{
				ToolID:    toolID,
				FileSet:   fileSet,
				Arguments: append([]string{}, arguments...),
			},
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
	execution := rule.Spec.Detail.(contract.FileCommandExecution)
	execution.ConfigArgument = configArgument
	execution.ConfigFile = configFile
	rule.Spec.Detail = execution
	return rule
}

func golangciRule(
	id string,
	name string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:    id,
		Name:  name,
		Group: RuleGroupLanguage,
		Spec: contract.ExecutionSpec{
			Kind: ExecutorBackendCommand,
			Detail: contract.BackendCommandExecution{
				ToolIDs: []string{
					ToolGo,
					ToolGoimports,
					ToolGolangciLint,
				},
				Action:   BackendActionGolangci,
				Language: LanguageGo,
			},
		},
		FixSpec: contract.ExecutionSpec{
			Kind: ExecutorBackendCommand,
			Detail: contract.BackendCommandExecution{
				ToolIDs: []string{
					ToolGo,
					ToolGoimports,
				},
				Action:   BackendActionGoFormat,
				Language: LanguageGo,
			},
		},
	}
}

func goStyleRule(
	id string,
	name string,
	check string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:    id,
		Name:  name,
		Group: RuleGroupLanguage,
		Spec: contract.ExecutionSpec{
			Kind: ExecutorBackendCheck,
			Detail: contract.BackendCheckExecution{
				ToolIDs:  []string{ToolGo},
				Check:    check,
				Language: LanguageGo,
			},
		},
	}
}

func scannerRule(
	id string,
	name string,
	scanner string,
) (rule RuleDefinition) {
	return scanRuleWithConfig(id, name, RuleGroupText, scanner)
}

func scanRuleWithConfig(
	id string,
	name string,
	group contract.RuleGroup,
	scanner string,
	configRefs ...string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:                 id,
		Name:               name,
		Group:              group,
		RequiredConfigRefs: append([]string{}, configRefs...),
		Spec: contract.ExecutionSpec{
			Kind: ExecutorRepositoryScan,
			Detail: contract.RepositoryScanExecution{
				Scanner: scanner,
			},
		},
	}
}

func lineLengthRule() (rule RuleDefinition) {
	rule = scannerRule(
		"text/line-length",
		"Line length",
		ScannerLineLength,
	)
	execution := rule.Spec.Detail.(contract.RepositoryScanExecution)
	execution.FileSet = "line_length"
	rule.Spec.Detail = execution
	return rule
}
