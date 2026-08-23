package profile_test

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
)

func TestValidateChecksProfile(t *testing.T) {
	t.Parallel()

	err := profile.Validate(profile.Profile{SchemaVersion: 2})
	requireErrorContains(t, err, "version 2")
}
