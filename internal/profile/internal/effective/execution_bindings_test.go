package effective_test

import (
	"testing"

	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/profiletest"
)

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.FileSets = nil
	_, err := effective.Compile(config, profiletest.FileCommandDefinitions())
	requireErrorContains(t, err, "unknown file set")
}
