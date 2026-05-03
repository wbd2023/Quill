package policy_test

import (
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
)

func TestLanguageConfigLookupBackend(t *testing.T) {
	language := policy.LanguageConfig{
		Backends: []policy.LanguageBackendConfig{
			{Name: "tooling_go", Language: "go", Scope: "tools"},
		},
	}

	backend, found := language.LookupBackend("tooling_go")
	if !found {
		t.Fatalf("expected language backend lookup to find tooling_go")
	}

	requireEqual(t, policy.LanguageBackendConfig{
		Name:     "tooling_go",
		Language: "go",
		Scope:    contract.Scope("tools"),
	}, backend)

	_, found = language.LookupBackend("missing")
	if found {
		t.Fatalf("expected missing language backend lookup to fail")
	}
}
