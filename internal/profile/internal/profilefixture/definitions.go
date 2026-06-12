package profilefixture

import "ciphera/tools/internal/style"

// Definitions returns rule and tool definitions that match Config.
func Definitions() (definitions style.Definitions) {
	return style.Definitions{
		Tools: []style.Tool{
			{ID: Tool, Name: "Test tool"},
		},
		Rules: []style.RuleDefinition{
			{
				ID:    Rule,
				Name:  "Test rule",
				Group: "test",
				Check: style.ExecutionSpec{
					Kind: style.ExecutionRepositoryScan,
					Detail: style.RepositoryScanExecution{
						Scanner: "test",
					},
				},
			},
		},
	}
}

// FileCommandDefinitions returns definitions with a file-command rule.
func FileCommandDefinitions() (definitions style.Definitions) {
	definitions = Definitions()
	definitions.Rules[0].Check = style.ExecutionSpec{
		Kind: style.ExecutionFileCommand,
		Detail: style.FileCommandExecution{
			ToolID:  Tool,
			FileSet: FileSet,
		},
	}
	return definitions
}

// TargetCommandDefinitions returns definitions with target check and fix executions.
func TargetCommandDefinitions() (definitions style.Definitions) {
	definitions = Definitions()
	definitions.Rules[0].Check = style.ExecutionSpec{
		Kind: style.ExecutionTargetCommand,
		Detail: style.TargetCommandExecution{
			ToolIDs:  []string{Tool},
			Action:   "test",
			Language: Language,
		},
	}
	definitions.Rules[0].Fix = style.ExecutionSpec{
		Kind: style.ExecutionTargetCommand,
		Detail: style.TargetCommandExecution{
			ToolIDs:  []string{Tool},
			Action:   "fix",
			Language: Language,
		},
	}
	return definitions
}

// TargetCheckDefinitions returns definitions with a target check execution.
func TargetCheckDefinitions() (definitions style.Definitions) {
	definitions = Definitions()
	definitions.Rules[0].Check = style.ExecutionSpec{
		Kind: style.ExecutionTargetCheck,
		Detail: style.TargetCheckExecution{
			ToolIDs:  []string{Tool},
			Check:    "test",
			Language: Language,
		},
	}
	return definitions
}
