package builtin

func securityPack() (pack Pack) {
	return Pack{
		ID:   PackSecurity,
		Name: "Security",
		Rules: []RuleDefinition{
			scanRule(
				"security/secrets",
				"Committed secrets",
				RuleGroupSecurity,
				ScannerSecrets,
			),
		},
	}
}
