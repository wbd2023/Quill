package profiles

import (
	"testing"

	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/testutil"
)

// Self loads this repository's Profile for explicit self-hosting tests.
func Self(test *testing.T) (config profile.Profile) {
	test.Helper()

	config, err := profile.Load(testutil.RepositoryRoot(test))
	if err != nil {
		test.Fatalf("profile.Load: %v", err)
	}

	return config
}

// SelfRepositoryConfig returns this repository's collector policy.
func SelfRepositoryConfig(test *testing.T) (repository profile.RepositoryConfig) {
	test.Helper()

	return Self(test).Repository
}

// Write writes the profile and STYLE.md to the given root.
func Write(test *testing.T, root string, config profile.Profile) {
	test.Helper()

	styleGuide := testutil.ReadFile(test, testutil.RepositoryRoot(test), "STYLE.md")
	testutil.WriteFile(test, root, config.StyleGuide.Path, styleGuide)
	testutil.WriteFile(test, root, "quill.toml", Format(test, config))
}

// Format serialises a profile config to its TOML representation.
func Format(test *testing.T, config profile.Profile) (contents string) {
	test.Helper()

	contents, err := profile.Format(config)
	if err != nil {
		test.Fatalf("format profile TOML: %v", err)
	}

	return contents
}
