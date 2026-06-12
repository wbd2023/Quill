package effective_test

import (
	"testing"

	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/profilefixture"
)

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	t.Parallel()

	config := profilefixture.Config()

	config.FileSets = nil
	_, err := effective.Compile(config, profilefixture.FileCommandDefinitions())
	requireErrorContains(t, err, "unknown file set")
}
