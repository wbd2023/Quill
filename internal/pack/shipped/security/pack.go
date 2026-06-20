package security

import (
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/style"
)

// PackID is pack i d.
const PackID = "security"

// ScannerSecrets is scanner secrets.
const ScannerSecrets = "secrets"

const ruleGroupSecurity style.RuleGroup = "security_scanners"

// Pack returns the Security Shipped Pack definition.
func Pack() (definition pack.Definition) {
	return pack.Definition{
		ID:   PackID,
		Name: "Security",
		Rules: []style.RuleDefinition{
			{
				ID:    "security/secrets",
				Name:  "Committed secrets",
				Group: ruleGroupSecurity,
				Check: style.ExecutionSpec{
					Kind: style.ExecutionRepositoryScan,
					Detail: style.RepositoryScanExecution{
						Scanner: ScannerSecrets,
					},
				},
			},
		},
	}
}
