package rulepack

func goPack() (pack Pack) {
	return Pack{
		ID:   PackGo,
		Name: "Go",
		Tools: selectTools(
			ToolGo,
			ToolGoimports,
			ToolGolangciLint,
		),
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
		goStyleRule("go/comments", "Go comments", GoCheckComments),
		goStyleRule("go/errors", "Go error handling", GoCheckErrors),
		goStyleRule("go/resources", "Go resource handling", GoCheckResources),
		goStyleRule("go/returns", "Go return style", GoCheckReturns),
		goStyleRule("go/parameters", "Go parameter style", GoCheckParameters),
		goStyleRule("go/domain-identifiers", "Go domain identifiers", GoCheckDomainIdentifiers),
		goStyleRule("go/naming", "Go naming", GoCheckNaming),
		goStyleRule("go/order", "Go declaration and method order", GoCheckOrder),
		goStyleRule("go/logging", "Go logging", GoCheckLogging),
		goStyleRule("go/security", "Go security", GoCheckSecurity),
		goStyleRule("go/process", "Go process execution", GoCheckProcess),
		goStyleRule("go/data", "Go data usage", GoCheckData),
		goStyleRule("go/tests", "Go test hygiene", GoCheckTests),
		goStyleRule("go/file-shape", "Go file shape", GoCheckFileShape),
	}
}

func goSpacingRules() (rules []RuleDefinition) {
	return []RuleDefinition{
		goStyleRule(
			"go/guard-clause-spacing",
			"Guard-clause spacing (Go)",
			GoCheckGuardClauseSpacing,
		),
		goStyleRule(
			"go/switch-case-spacing",
			"Switch-case spacing (Go)",
			GoCheckSwitchCaseSpacing,
		),
	}
}

func goArchitectureRules() (rules []RuleDefinition) {
	return []RuleDefinition{
		scanRuleWithConfig(
			"go/architecture-imports",
			"Architecture imports",
			RuleGroupLanguage,
			ScannerArchitecture,
			ConfigReferenceGo,
		),
	}
}
