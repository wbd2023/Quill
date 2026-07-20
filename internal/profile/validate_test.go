package profile_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/profile"
)

func TestValidateChecksProfile(t *testing.T) {
	t.Parallel()

	err := profile.Validate(policy.Config{SchemaVersion: 2})
	requireErrorContains(t, err, "version 2")
}
