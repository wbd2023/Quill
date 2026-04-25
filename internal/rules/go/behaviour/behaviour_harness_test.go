package behaviour

import (
	"testing"

	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
	"ciphera/tools/internal/profile"
	gostyle "ciphera/tools/internal/rules/go"
)

/* ------------------------------------------- Harness ------------------------------------------ */

func runGoStyleCheck(t *testing.T, targetDirectory string) (output string, err error) {
	t.Helper()

	return runGoStyleCheckWithPolicy(t, targetDirectory, profiles.Current(t))
}

func runGoStyleCheckWithPolicy(
	t *testing.T,
	targetDirectory string,
	policy profile.Profile,
) (output string, err error) {
	t.Helper()

	return gostyle.CheckDirectories(
		targetDirectory,
		[]string{targetDirectory},
		policy,
	)
}

func writeTypeAwareDomainFixture(t *testing.T, rootDirectory string) {
	t.Helper()

	fixtures.WriteFile(t, rootDirectory, "go.mod", "module example\n\ngo 1.24.5\n")
	fixtures.WriteFile(
		t,
		rootDirectory,
		"internal/core/domain/types.go",
		`package domain

type IdentityID string
`,
	)
}

func writeSourceFile(t *testing.T, path string, contents string) {
	t.Helper()

	fixtures.WritePath(t, path, contents)
}
