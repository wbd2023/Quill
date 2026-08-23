package golang

import (
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/style"
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
		styleRule("go/comments", "Go comments", TargetCheckComments),
		styleRule("go/errors", "Go error handling", TargetCheckErrors),
		styleRule("go/resources", "Go resource handling", TargetCheckResources),
		styleRule("go/returns", "Go return style", TargetCheckReturns),
		styleRule("go/parameters", "Go parameter style", TargetCheckParameters),
		styleRule("go/domain-values", "Go domain values", TargetCheckDomainValues),
		styleRule("go/naming", "Go naming", TargetCheckNaming),
		styleRule("go/order", "Go declaration and method order", TargetCheckOrder),
		styleRule("go/logging", "Go logging", TargetCheckLogging),
		styleRule("go/security", "Go security", TargetCheckSecurity),
		styleRule("go/process", "Go process execution", TargetCheckProcess),
		styleRule("go/data", "Go data usage", TargetCheckData),
		styleRule("go/tests", "Go test hygiene", TargetCheckTests),
		styleRule("go/file-shape", "Go file shape", TargetCheckFileShape),
	}
}

func spacingRules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{
		styleRule(
			"go/guard-clause-spacing",
			"Guard-clause spacing (Go)",
			TargetCheckGuardClauseSpacing,
		),
		styleRule(
			"go/switch-case-spacing",
			"Switch-case spacing (Go)",
			TargetCheckSwitchCaseSpacing,
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
		Check: style.TargetCommandTemplate{
			ToolIDs: []string{
				tool.Go,
				tool.Goimports,
				tool.GolangciLint,
			},
			Action:   TargetActionGolangci,
			Language: Language,
		},
		Fix: style.TargetCommandTemplate{
			ToolIDs: []string{
				tool.Go,
				tool.Goimports,
			},
			Action:   TargetActionGoFormat,
			Language: Language,
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
		Check: style.TargetCheckTemplate{
			ToolIDs:  []string{tool.Go},
			Check:    checkID,
			Language: Language,
		},
	}
}

func architectureRule() (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    "go/architecture-imports",
		Name:  "Architecture imports",
		Group: ruleGroupLanguage,
		Check: style.RepositoryScan{
			Scanner: ScannerArchitecture,
		},
	}
}
