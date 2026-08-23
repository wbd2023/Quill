package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
)

func TestCompileRequiresActivePinnedTools(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Tools = config.Tools[1:]

	_, err := profile.Compile(config, profiletest.Definitions())
	requireErrorContains(t, err, "missing a pinned tool")
}

func TestCompileRejectsUnknownPinnedTools(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.Tools = append(config.Tools, profile.PinnedTool{
		ID:      "unknown",
		Version: "1.0.0",
	})
	_, err := profile.Compile(config, profiletest.Definitions())
	requireErrorContains(t, err, "unknown")
}
