package effective_test

import (
	"testing"

	"ciphera/tools/internal/profile/internal/effective"
	"ciphera/tools/internal/profile/internal/fixture"
)

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	t.Parallel()

	config := fixture.Config()

	config.FileSets = nil
	_, err := effective.Compile(config, fixture.FileCommandDefinitions())
	requireErrorContains(t, err, "unknown file set")
}
