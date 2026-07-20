package profiletest

import "github.com/wbd2023/Quill/internal/style"

// Definitions returns rule and tool definitions that match Config.
func Definitions() (definitions style.Definitions) {
	return style.Definitions{
		ToolIDs: []string{Tool},
		Rules: []style.RuleDefinition{
			{
				ID:    Rule,
				Name:  "Test rule",
				Group: "test",
				Check: style.RepositoryScanExecution{
					Scanner: "test",
				},
			},
		},
	}
}

// FileCommandDefinitions returns definitions with a file-command rule.
func FileCommandDefinitions() (definitions style.Definitions) {
	definitions = Definitions()
	definitions.Rules[0].Check = style.FileCommandExecution{
		ToolID:  Tool,
		FileSet: FileSet,
	}
	return definitions
}

// TargetCommandDefinitions returns definitions with target check and fix executions.
func TargetCommandDefinitions() (definitions style.Definitions) {
	definitions = Definitions()
	definitions.Rules[0].Check = style.TargetCommandTemplate{
		ToolIDs:  []string{Tool},
		Action:   "test",
		Language: Language,
	}
	definitions.Rules[0].Fix = style.TargetCommandTemplate{
		ToolIDs:  []string{Tool},
		Action:   "fix",
		Language: Language,
	}
	return definitions
}

// TargetCheckDefinitions returns definitions with a target check execution.
func TargetCheckDefinitions() (definitions style.Definitions) {
	definitions = Definitions()
	definitions.Rules[0].Check = style.TargetCheckTemplate{
		ToolIDs:  []string{Tool},
		Check:    "test",
		Language: Language,
	}
	return definitions
}
