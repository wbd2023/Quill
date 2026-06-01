package profile_test

import (
	"testing"

	"ciphera/tools/internal/profile"
)

func TestParseValidatesProfile(t *testing.T) {
	t.Parallel()

	_, err := profile.Parse(`schema_version = 2`)
	requireErrorContains(t, err, "version 2")
}
