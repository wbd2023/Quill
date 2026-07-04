package profiles

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/testutil"
)

// Current loads the active profile from the repository root.
func Current(test *testing.T) (config policy.Config) {
	test.Helper()

	config, err := profile.Load(testutil.RepositoryRoot(test))
	if err != nil {
		test.Fatalf("profile.Load: %v", err)
	}

	return config
}

// RepositoryConfig repository config.
func RepositoryConfig(test *testing.T) (repository policy.RepositoryConfig) {
	test.Helper()

	return Current(test).Repository
}

// Write writes the profile and STYLE.md to the given root.
func Write(test *testing.T, root string, config policy.Config) {
	test.Helper()

	styleGuide := testutil.ReadFile(test, testutil.RepositoryRoot(test), "STYLE.md")
	testutil.WriteFile(test, root, config.StyleGuide.Path, styleGuide)
	testutil.WriteFile(test, root, "style.toml", Format(test, config))
}

// Format serialises a profile config to its TOML representation.
func Format(test *testing.T, config policy.Config) (contents string) {
	test.Helper()

	contents, err := profile.Format(config)
	if err != nil {
		test.Fatalf("format profile TOML: %v", err)
	}

	return contents
}
