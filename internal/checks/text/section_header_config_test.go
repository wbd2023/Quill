package text

import (
	"testing"

	"github.com/wbd2023/Quill/internal/checks/textpolicy"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func currentSectionHeaders(t *testing.T) (headers textpolicy.SectionHeaderConfig) {
	t.Helper()

	pack, found := profiles.Current(t).PackConfigs.Lookup("text")
	if !found {
		t.Fatal("expected text pack config")
	}

	config, err := textpolicy.DecodeConfig(pack)
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	return config.SectionHeaders
}
