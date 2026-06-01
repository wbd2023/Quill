package builtin

import "ciphera/tools/internal/contract"

const (
	RuleGroupProject    contract.RuleGroup = "project"
	RuleGroupLanguage   contract.RuleGroup = "language"
	RuleGroupText       contract.RuleGroup = "text_scanners"
	RuleGroupSecurity   contract.RuleGroup = "security_scanners"
	RuleGroupVocabulary contract.RuleGroup = "vocabulary_scanners"
	RuleGroupExternal   contract.RuleGroup = "external_tools"
)
