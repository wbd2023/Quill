package rulepack

/* ------------------------------------------- Go Pack ------------------------------------------ */

func goPack() (pack Pack) {
	return Pack{
		ID:   PackGo,
		Name: "Go",
		Tools: selectTools(
			ToolGo,
			ToolGoimports,
			ToolGolangciLint,
		),
		Rules: []RuleDefinition{
			golangciRule(
				"go/lint",
				"golangci-lint",
			),
			goStyleRule(
				"go/comments",
				"Go comments",
				GoCheckComments,
			),
			goStyleRule(
				"go/errors",
				"Go error handling",
				GoCheckErrors,
			),
			goStyleRule(
				"go/resources",
				"Go resource handling",
				GoCheckResources,
			),
			goStyleRule(
				"go/returns",
				"Go return style",
				GoCheckReturns,
			),
			goStyleRule(
				"go/parameters",
				"Go parameter style",
				GoCheckParameters,
			),
			goStyleRule(
				"go/domain-identifiers",
				"Go domain identifiers",
				GoCheckDomainIdentifiers,
			),
			goStyleRule(
				"go/naming",
				"Go naming",
				GoCheckNaming,
			),
			goStyleRule(
				"go/order",
				"Go declaration and method order",
				GoCheckOrder,
			),
			goStyleRule(
				"go/logging",
				"Go logging",
				GoCheckLogging,
			),
			goStyleRule(
				"go/security",
				"Go security",
				GoCheckSecurity,
			),
			goStyleRule(
				"go/process",
				"Go process execution",
				GoCheckProcess,
			),
			goStyleRule(
				"go/data",
				"Go data usage",
				GoCheckData,
			),
			goStyleRule(
				"go/tests",
				"Go test hygiene",
				GoCheckTests,
			),
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
			scanRuleWithConfig(
				"go/architecture-imports",
				"Architecture imports",
				RuleGroupLanguage,
				ScannerArchitecture,
				ConfigRefArchitecture,
			),
		},
	}
}
