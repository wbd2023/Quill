package text

import (
	"testing"

	textpolicy "github.com/wbd2023/quill/internal/pack/shipped/text/policy"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

func currentSectionHeaders(t *testing.T) (headers textpolicy.SectionHeaderConfig) {
	t.Helper()

	pack, found := profiles.Self(t).PackConfigs.Lookup("text")
	if !found {
		t.Fatal("expected text pack config")
	}

	config, err := textpolicy.DecodeConfig(pack)
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	return config.SectionHeaders
}
