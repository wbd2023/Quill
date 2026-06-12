package profiles

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/testutil"
)

func Current(test *testing.T) (config policy.Config) {
	test.Helper()

	config, err := profile.Load(testutil.RepositoryRoot(test))
	if err != nil {
		test.Fatalf("profile.Load: %v", err)
	}

	return config
}

func RepositoryConfig(test *testing.T) (repository policy.RepositoryConfig) {
	test.Helper()

	return Current(test).Repository
}

func Write(test *testing.T, root string, config policy.Config) {
	test.Helper()

	styleGuide := testutil.ReadFile(test, testutil.RepositoryRoot(test), "STYLE.md")
	testutil.WriteFile(test, root, config.StyleGuide.Path, styleGuide)
	testutil.WriteFile(test, root, "style.toml", Format(test, config))
}

func Format(test *testing.T, config policy.Config) (contents string) {
	test.Helper()

	contents, err := profile.Format(config)
	if err != nil {
		test.Fatalf("format profile TOML: %v", err)
	}

	return contents
}
