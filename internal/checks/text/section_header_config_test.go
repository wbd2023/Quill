package text

import (
	"testing"

	textpolicy "github.com/wbd2023/quill/internal/pack/shipped/text/policy"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

func currentSectionHeaders(t *testing.T) (headers textpolicy.SectionHeaderConfig) {
	t.Helper()

	policy, found := profiles.Self(t).PackPolicies.Lookup("text")
	if !found {
		t.Fatal("expected Text Pack Policy")
	}

	config, err := textpolicy.DecodeConfig(policy)
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	return config.SectionHeaders
}
