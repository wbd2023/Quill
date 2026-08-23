package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/profile/internal/profiletest"
)

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	t.Parallel()

	config := profiletest.Config()
	config.FileSets = nil

	_, err := profile.Compile(config, profiletest.FileCommandDefinitions())
	requireErrorContains(t, err, "unknown file set")
}
