package profile

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func effectiveRuleByID(
	effective contract.EffectiveConfig,
	ruleID string,
) (rule contract.Rule, found bool) {
	for _, rule := range effective.Rules {
		if rule.ID == ruleID {
			return rule, true
		}
	}

	return contract.Rule{}, false
}

func replaceRuleBindingBackends(
	config *policy.Config,
	ruleID string,
	backends []string,
) {
	for index := range config.Rules {
		if config.Rules[index].RuleID != ruleID {
			continue
		}

		config.Rules[index].Backends = append([]string{}, backends...)
		return
	}
}
