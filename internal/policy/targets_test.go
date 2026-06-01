package policy_test

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func TestTargetConfigsLookup(t *testing.T) {
	var targets policy.TargetConfigs
	targets = append(targets, policy.TargetConfig{
		Name:     "tools_go",
		Language: "go",
		Scope:    "tools",
	})

	target, found := targets.Lookup("tools_go")
	if !found {
		t.Fatalf("expected target lookup to find tools_go")
	}

	requireEqual(t, policy.TargetConfig{
		Name:     "tools_go",
		Language: "go",
		Scope:    contract.Scope("tools"),
	}, target)

	_, found = targets.Lookup("missing")
	if found {
		t.Fatalf("expected missing target lookup to fail")
	}
}
