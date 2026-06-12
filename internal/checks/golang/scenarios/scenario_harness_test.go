package scenarios

import (
	"testing"

	"ciphera/tools/internal/checks/golang"
	"ciphera/tools/internal/checks/gopolicy"
	gopack "ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

func runGoStyleResult(
	t *testing.T,
	targetDirectory string,
) (result style.ExecutionResult, err error) {
	t.Helper()

	return runGoStyleResultWithPolicy(t, targetDirectory, profiles.Current(t))
}

func runGoStyleResultWithPolicy(
	t *testing.T,
	targetDirectory string,
	config policy.Config,
) (result style.ExecutionResult, err error) {
	t.Helper()

	result, err = golang.CheckDirectories(
		targetDirectory,
		[]string{targetDirectory},
		config.Repository,
		config.PathRoles,
		goConfigForTest(t, config),
	)
	return result, err
}

func goConfigForTest(t *testing.T, config policy.Config) (goConfig gopolicy.Config) {
	t.Helper()

	pack, found := config.PackConfigs.Lookup(gopack.PackID)
	if !found {
		t.Fatal("missing Go pack config")
	}

	goConfig, err := gopolicy.DecodeConfig(pack)
	if err != nil {
		t.Fatalf("Decode Go config: %v", err)
	}

	return goConfig
}

func updateGoConfigForTest(
	t *testing.T,
	config *policy.Config,
	update func(*gopolicy.Config),
) {
	t.Helper()

	goConfig := goConfigForTest(t, *config)
	update(&goConfig)
	config.PackConfigs[gopack.PackID] = gopolicy.EncodeConfig(goConfig)
}

func writeTypeAwareDomainFixture(t *testing.T, rootDirectory string) {
	t.Helper()

	testutil.WriteFile(t, rootDirectory, "go.mod", "module example\n\ngo 1.24.5\n")
	testutil.WriteFile(
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

	testutil.WritePath(t, path, contents)
}
