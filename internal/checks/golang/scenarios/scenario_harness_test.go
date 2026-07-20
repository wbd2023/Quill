package scenarios

import (
	"testing"

	"github.com/wbd2023/Quill/internal/checks/golang"
	"github.com/wbd2023/Quill/internal/checks/gopolicy"
	gopack "github.com/wbd2023/Quill/internal/pack/shipped/golang"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/testutil"
	"github.com/wbd2023/Quill/internal/testutil/profiles"
)

func runGoStyleResult(
	t *testing.T,
	targetDirectory string,
) (result style.ExecutionResult, err error) {
	t.Helper()

	return runGoStyleResultWithPolicy(t, targetDirectory, scenarioConfig(t))
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

/* -------------------------------------- Scenario Profile -------------------------------------- */

func scenarioConfig(t *testing.T) (config policy.Config) {
	t.Helper()

	config = profiles.Current(t)
	config.PathRoles = policy.PathRoles{
		"go_source": {"cmd/", "internal/", "test/"},
		"application_port": {
			"internal/client/application/port/",
			"internal/relay/application/port/",
		},
		"concrete_infra": {"internal/client/adapters/", "internal/relay/adapters/"},
		"domain":         {"internal/core/domain/"},
		"domain_errors":  {"internal/core/domain/errors.go"},
		"test_mocks":     {"internal/testkit/mocks/"},
	}

	goConfig := gopolicy.Config{
		LocalImportPrefixes: []string{"ciphera"},
		Parameters: gopolicy.ParameterConfig{
			SecretNames: []string{
				"passphrase",
				"privateKey",
				"token",
				"seed",
				"secret",
				"password",
				"secretKey",
			},
		},
		Constructors: gopolicy.ConstructorConfig{
			ParameterOrder: []gopolicy.ParameterGroup{
				{Name: "repository", TypeNameSuffixes: []string{"Repository"}},
				{Name: "service", TypeNameSuffixes: []string{"Service"}},
				{Name: "adapter", TypeNameSuffixes: []string{"Client", "Factory"}},
				{
					Name:           "config",
					ParameterNames: []string{"serverURL", "relayURL", "identityID", "timeout"},
				},
				{Name: "secret", MatchesSecretNames: true},
			},
		},
		DomainValues: gopolicy.DomainValueConfig{
			RequiredConstructors: gopolicy.DomainValueConstructors{
				"Username":       {"ParseUsername"},
				"ConversationID": {"ParseConversationID", "ConversationIDFromUsername"},
				"IdentityID":     {"ParseIdentityID"},
			},
		},
		Architecture: gopolicy.ArchitectureConfig{
			Layers: []gopolicy.ArchitectureLayer{
				{
					Name:          "core",
					PackageRoots:  []string{"internal/core"},
					AllowedLayers: []string{"core"},
				},
				{
					Name:          "client_port",
					PackageRoots:  []string{"internal/client/application/port"},
					AllowedLayers: []string{"core", "client_port"},
				},
				{
					Name:          "client_service",
					PackageRoots:  []string{"internal/client/application/service"},
					AllowedLayers: []string{"core", "client_port", "client_service"},
				},
				{
					Name:         "client_inbound",
					PackageRoots: []string{"internal/client/adapters/inbound"},
					AllowedLayers: []string{
						"core",
						"client_port",
						"client_service",
						"client_inbound",
						"client_bootstrap",
						"shared",
					},
				},
				{
					Name:          "client_outbound",
					PackageRoots:  []string{"internal/client/adapters/outbound"},
					AllowedLayers: []string{"core", "client_port", "client_outbound", "shared"},
				},
				{
					Name:         "client_bootstrap",
					PackageRoots: []string{"internal/client/bootstrap"},
					AllowedLayers: []string{
						"core",
						"client_port",
						"client_service",
						"client_inbound",
						"client_outbound",
						"client_bootstrap",
						"shared",
					},
				},
				{
					Name:          "shared",
					PackageRoots:  []string{"internal/relaywire"},
					AllowedLayers: []string{"core", "shared"},
				},
			},
		},
	}
	config.PackConfigs[gopack.PackID] = gopolicy.EncodeConfig(goConfig)

	return config
}

/* --------------------------------------- Config Updates --------------------------------------- */

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

	testutil.WriteFile(t, "", path, contents)
}
