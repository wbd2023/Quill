package profile

import (
	"testing"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/profile/internal/profiletest"
)

func TestCompileRequiresActivePinnedTools(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.Tools = config.Tools[1:]
	_, err := compilePlan(config, profiletest.Definitions())
	requireErrorContainsInternal(t, err, "missing a pinned tool")
}

func TestCompileRejectsUnknownPinnedTools(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.Tools = append(config.Tools, policy.PinnedTool{ID: "unknown", Version: "1.0.0"})
	_, err := compilePlan(config, profiletest.Definitions())
	requireErrorContainsInternal(t, err, "unknown")
}
