package builtin

import (
	"ciphera/tools/internal/rules/golang/check"
	gopolicy "ciphera/tools/internal/rules/golang/policy"
)

func goPack() (pack Pack) {
	return Pack{
		ID:   PackGo,
		Name: "Go",
		Tools: selectTools(
			ToolGo,
			ToolGoimports,
			ToolGolangciLint,
		),
		Config: PackConfig{
			Required: true,
			Validate: gopolicy.ValidatePackConfig,
		},
		Rules: goRules(),
	}
}

func goRules() (rules []RuleDefinition) {
	rules = append(rules, goToolRules()...)
	rules = append(rules, goStructuredRules()...)
	rules = append(rules, goSpacingRules()...)
	rules = append(rules, goArchitectureRules()...)
	return rules
}

func goToolRules() (rules []RuleDefinition) {
	return []RuleDefinition{
		golangciRule(
			"go/lint",
			"golangci-lint",
		),
	}
}

func goStructuredRules() (rules []RuleDefinition) {
	return []RuleDefinition{
		goStyleRule("go/comments", "Go comments", check.Comments),
		goStyleRule("go/errors", "Go error handling", check.Errors),
		goStyleRule("go/resources", "Go resource handling", check.Resources),
		goStyleRule("go/returns", "Go return style", check.Returns),
		goStyleRule("go/parameters", "Go parameter style", check.Parameters),
		goStyleRule("go/domain-values", "Go domain values", check.DomainValues),
		goStyleRule("go/naming", "Go naming", check.Naming),
		goStyleRule("go/order", "Go declaration and method order", check.Order),
		goStyleRule("go/logging", "Go logging", check.Logging),
		goStyleRule("go/security", "Go security", check.Security),
		goStyleRule("go/process", "Go process execution", check.Process),
		goStyleRule("go/data", "Go data usage", check.Data),
		goStyleRule("go/tests", "Go test hygiene", check.Tests),
		goStyleRule("go/file-shape", "Go file shape", check.FileShape),
	}
}

func goSpacingRules() (rules []RuleDefinition) {
	return []RuleDefinition{
		goStyleRule(
			"go/guard-clause-spacing",
			"Guard-clause spacing (Go)",
			check.GuardClauseSpacing,
		),
		goStyleRule(
			"go/switch-case-spacing",
			"Switch-case spacing (Go)",
			check.SwitchCaseSpacing,
		),
	}
}

func goArchitectureRules() (rules []RuleDefinition) {
	return []RuleDefinition{
		scanRule(
			"go/architecture-imports",
			"Architecture imports",
			RuleGroupLanguage,
			ScannerArchitecture,
		),
	}
}
