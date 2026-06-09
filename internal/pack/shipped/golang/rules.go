package golang

import (
	"ciphera/tools/internal/checks/golang/check"
	"ciphera/tools/internal/style"
)

const (
	TargetActionGoFormat = "go_format"
	TargetActionGolangci = "golangci"

	Language            = "go"
	ScannerArchitecture = "architecture"
)

const ruleGroupLanguage style.RuleGroup = "language"

/* ----------------------------------------- Rule Lists ----------------------------------------- */

func rules() (rules []style.RuleDefinition) {
	rules = append(rules, toolRules()...)
	rules = append(rules, structuredRules()...)
	rules = append(rules, spacingRules()...)
	rules = append(rules, architectureRules()...)
	return rules
}

func toolRules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{
		golangciRule(
			"go/lint",
			"golangci-lint",
		),
	}
}

func structuredRules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{
		styleRule("go/comments", "Go comments", check.Comments),
		styleRule("go/errors", "Go error handling", check.Errors),
		styleRule("go/resources", "Go resource handling", check.Resources),
		styleRule("go/returns", "Go return style", check.Returns),
		styleRule("go/parameters", "Go parameter style", check.Parameters),
		styleRule("go/domain-values", "Go domain values", check.DomainValues),
		styleRule("go/naming", "Go naming", check.Naming),
		styleRule("go/order", "Go declaration and method order", check.Order),
		styleRule("go/logging", "Go logging", check.Logging),
		styleRule("go/security", "Go security", check.Security),
		styleRule("go/process", "Go process execution", check.Process),
		styleRule("go/data", "Go data usage", check.Data),
		styleRule("go/tests", "Go test hygiene", check.Tests),
		styleRule("go/file-shape", "Go file shape", check.FileShape),
	}
}

func spacingRules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{
		styleRule(
			"go/guard-clause-spacing",
			"Guard-clause spacing (Go)",
			check.GuardClauseSpacing,
		),
		styleRule(
			"go/switch-case-spacing",
			"Switch-case spacing (Go)",
			check.SwitchCaseSpacing,
		),
	}
}

func architectureRules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{architectureRule()}
}

/* ---------------------------------------- Rule Builders --------------------------------------- */

func golangciRule(
	id string,
	name string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupLanguage,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionTargetCommand,
			Detail: style.TargetCommandExecution{
				ToolIDs: []string{
					ToolGo,
					ToolGoimports,
					ToolGolangciLint,
				},
				Action:   TargetActionGolangci,
				Language: Language,
			},
		},
		Fix: style.ExecutionSpec{
			Kind: style.ExecutionTargetCommand,
			Detail: style.TargetCommandExecution{
				ToolIDs: []string{
					ToolGo,
					ToolGoimports,
				},
				Action:   TargetActionGoFormat,
				Language: Language,
			},
		},
	}
}

func styleRule(
	id string,
	name string,
	checkID string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupLanguage,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionTargetCheck,
			Detail: style.TargetCheckExecution{
				ToolIDs:  []string{ToolGo},
				Check:    checkID,
				Language: Language,
			},
		},
	}
}

func architectureRule() (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    "go/architecture-imports",
		Name:  "Architecture imports",
		Group: ruleGroupLanguage,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionRepositoryScan,
			Detail: style.RepositoryScanExecution{
				Scanner: ScannerArchitecture,
			},
		},
	}
}
