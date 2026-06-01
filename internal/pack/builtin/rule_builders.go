package builtin

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
		Group: RuleGroupProject,
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorToolchain,
			Detail: contract.ToolchainExecution{
				ToolIDs: append([]string{}, toolIDs...),
			},
		},
	}
}

func projectRule(
	id string,
	name string,
	check string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:    id,
		Name:  name,
		Group: RuleGroupProject,
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorProject,
			Detail: contract.ProjectExecution{
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
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorFileCommand,
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
	execution := rule.Check.Detail.(contract.FileCommandExecution)
	execution.ConfigArgument = configArgument
	execution.ConfigFile = configFile
	rule.Check.Detail = execution
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
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorTargetCommand,
			Detail: contract.TargetCommandExecution{
				ToolIDs: []string{
					ToolGo,
					ToolGoimports,
					ToolGolangciLint,
				},
				Action:   TargetActionGolangci,
				Language: LanguageGo,
			},
		},
		Fix: contract.ExecutionSpec{
			Kind: contract.ExecutorTargetCommand,
			Detail: contract.TargetCommandExecution{
				ToolIDs: []string{
					ToolGo,
					ToolGoimports,
				},
				Action:   TargetActionGoFormat,
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
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorTargetCheck,
			Detail: contract.TargetCheckExecution{
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
	return scanRule(id, name, RuleGroupText, scanner)
}

func scanRule(
	id string,
	name string,
	group contract.RuleGroup,
	scanner string,
) (rule RuleDefinition) {
	return RuleDefinition{
		ID:    id,
		Name:  name,
		Group: group,
		Check: contract.ExecutionSpec{
			Kind: contract.ExecutorRepositoryScan,
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
	execution := rule.Check.Detail.(contract.RepositoryScanExecution)
	execution.FileSet = "line_length"
	rule.Check.Detail = execution
	return rule
}
