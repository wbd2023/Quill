package effective_test

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/fixture"
)

func TestCompileRequiresActivePinnedTools(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.Tools = config.Tools[1:]
	_, err := effective.Compile(config, fixture.Definitions())
	requireErrorContains(t, err, "missing a pinned tool")
}

func TestCompileRejectsUnknownPinnedTools(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.Tools = append(config.Tools, policy.PinnedTool{ID: "unknown", Version: "1.0.0"})
	_, err := effective.Compile(config, fixture.Definitions())
	requireErrorContains(t, err, "unknown")
}
