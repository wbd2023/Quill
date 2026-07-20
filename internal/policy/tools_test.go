package policy_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/policy"
)

func TestPinnedToolsLookup(t *testing.T) {
	tools := policy.PinnedTools{
		{ID: "go", Version: "1.24.5", TimeoutSeconds: 30},
	}

	pinnedTool, found := tools.Lookup("go")
	if !found {
		t.Fatalf("expected pinned tool lookup to find go")
	}

	requireEqual(t, policy.PinnedTool{
		ID:             "go",
		Version:        "1.24.5",
		TimeoutSeconds: 30,
	}, pinnedTool)

	_, found = tools.Lookup("missing")
	if found {
		t.Fatalf("expected missing pinned tool lookup to fail")
	}
}
