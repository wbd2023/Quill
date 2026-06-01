package profiles

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile"
)

func Current(test *testing.T) (config policy.Config) {
	test.Helper()

	config, err := profile.Load(fixtures.RepositoryRoot(test))
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

	styleGuide := fixtures.ReadFile(test, fixtures.RepositoryRoot(test), "STYLE.md")
	fixtures.WriteFile(test, root, config.StyleGuide.Path, styleGuide)
	fixtures.WriteFile(test, root, "style.toml", Format(test, config))
}

func Format(test *testing.T, config policy.Config) (contents string) {
	test.Helper()

	contents, err := profile.Format(config)
	if err != nil {
		test.Fatalf("format profile TOML: %v", err)
	}

	return contents
}
