package engine

import (
	"context"
	"fmt"

	"github.com/wbd2023/quill/internal/profile"
)

// ExplainResult describes one active Rule and its owning Pack. It is presentation-free and is
// assembled without Tool inspection or subprocess execution.
type ExplainResult struct {
	Rule       RuleMetadata
	Pack       PackMetadata
	PackPolicy profile.PackPolicy
}

// Explain resolves one active rule from the repository's prepared metadata. A syntactically valid
// rule identifier that is unknown or inactive is an ArgumentError because it cannot name a rule in
// the selected profile.
func (engine *Engine) Explain(
	operationContext context.Context,
	ruleID string,
) (result ExplainResult, operationError error) {
	snapshot, err := engine.Metadata(operationContext)
	if err != nil {
		return ExplainResult{}, err
	}

	for _, rule := range snapshot.Rules {
		if rule.ID != ruleID {
			continue
		}
		if !rule.Active {
			return ExplainResult{}, newArgumentError(
				"rule %q is declared by pack %q but is not active in this profile",
				ruleID,
				rule.PackID,
			)
		}

		pack, found := findPack(rule.PackID, snapshot.Packs)
		if !found {
			return ExplainResult{}, fmt.Errorf(
				"rule %q references missing owner pack %q",
				rule.ID,
				rule.PackID,
			)
		}

		policy, _ := snapshot.PackPolicies.Lookup(rule.PackID)
		return ExplainResult{
			Rule:       rule,
			Pack:       pack,
			PackPolicy: policy.Clone(),
		}, nil
	}

	return ExplainResult{}, newArgumentError("unknown rule %q: no matching rule is declared", ruleID)
}

func findPack(packID string, packs []PackMetadata) (pack PackMetadata, found bool) {
	for _, candidate := range packs {
		if candidate.ID == packID {
			return candidate, true
		}
	}

	return PackMetadata{}, false
}
