package profile_test

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile"
)

func TestValidateChecksProfile(t *testing.T) {
	t.Parallel()

	err := profile.Validate(policy.Config{SchemaVersion: 2})
	requireErrorContains(t, err, "version 2")
}
