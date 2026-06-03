package security

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack"
)

const PackID = "security"

const ScannerSecrets = "secrets"

const ruleGroupSecurity contract.RuleGroup = "security_scanners"

// Pack returns the Security Shipped Pack definition.
func Pack() (definition pack.Definition) {
	return pack.Definition{
		ID:   PackID,
		Name: "Security",
		Rules: []contract.RuleDefinition{
			{
				ID:    "security/secrets",
				Name:  "Committed secrets",
				Group: ruleGroupSecurity,
				Check: contract.ExecutionSpec{
					Kind: contract.ExecutionRepositoryScan,
					Detail: contract.RepositoryScanExecution{
						Scanner: ScannerSecrets,
					},
				},
			},
		},
	}
}
