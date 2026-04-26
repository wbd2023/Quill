package rulepack

func securityPack() (pack Pack) {
	return Pack{
		ID:   PackSecurity,
		Name: "Security",
		Rules: []RuleDefinition{
			scanRuleWithConfig(
				"security/secrets",
				"Committed secrets",
				RuleGroupSecurity,
				ScannerSecrets,
			),
		},
	}
}
