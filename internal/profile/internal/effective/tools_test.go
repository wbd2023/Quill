package effective_test

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/profilefixture"
)

func TestCompileRequiresActivePinnedTools(t *testing.T) {
	t.Parallel()

	config := profilefixture.Config()

	config.Tools = config.Tools[1:]
	_, err := effective.Compile(config, profilefixture.Definitions())
	requireErrorContains(t, err, "missing a pinned tool")
}

func TestCompileRejectsUnknownPinnedTools(t *testing.T) {
	t.Parallel()

	config := profilefixture.Config()

	config.Tools = append(config.Tools, policy.PinnedTool{ID: "unknown", Version: "1.0.0"})
	_, err := effective.Compile(config, profilefixture.Definitions())
	requireErrorContains(t, err, "unknown")
}
