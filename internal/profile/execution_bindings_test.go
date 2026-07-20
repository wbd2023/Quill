package profile

import (
	"testing"

	"github.com/wbd2023/Quill/internal/profile/internal/profiletest"
)

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()

	config.FileSets = nil
	_, err := compilePlan(config, profiletest.FileCommandDefinitions())
	requireErrorContainsInternal(t, err, "unknown file set")
}
