package scenarios

import (
	"testing"

	"github.com/wbd2023/quill/internal/checks/golang"
	gopack "github.com/wbd2023/quill/internal/pack/shipped/golang"
	gopolicy "github.com/wbd2023/quill/internal/pack/shipped/golang/policy"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

func runGoStyleResult(
	t *testing.T,
	directory string,
) (result style.ExecutionResult, err error) {
	t.Helper()

	return runGoStyleResultWithPolicy(t, directory, scenarioConfig(t))
}

func runGoStyleResultWithPolicy(
	t *testing.T,
	directory string,
	config profile.Profile,
) (result style.ExecutionResult, err error) {
	t.Helper()

	result, err = golang.CheckDirectories(
		directory,
		[]string{directory},
		config.Repository,
		config.PathRoles,
		goConfigForTest(t, config),
	)
	return result, err
}

func goConfigForTest(t *testing.T, config profile.Profile) (goConfig gopolicy.Config) {
	t.Helper()

	policy, found := config.PackPolicies.Lookup(gopack.PackID)
	if !found {
		t.Fatal("missing Go Pack Policy")
	}

	goConfig, err := gopolicy.DecodeConfig(policy)
	if err != nil {
		t.Fatalf("Decode Go config: %v", err)
	}

	return goConfig
}

/* -------------------------------------- Scenario Profile -------------------------------------- */

func scenarioConfig(t *testing.T) (config profile.Profile) {
	t.Helper()

	config = profiles.Self(t)
	config.PathRoles = profile.PathRoles{
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
	config.PackPolicies[gopack.PackID] = gopolicy.EncodeConfig(goConfig)

	return config
}

/* --------------------------------------- Config Updates --------------------------------------- */

func updateGoConfigForTest(
	t *testing.T,
	config *profile.Profile,
	update func(*gopolicy.Config),
) {
	t.Helper()

	goConfig := goConfigForTest(t, *config)
	update(&goConfig)
	config.PackPolicies[gopack.PackID] = gopolicy.EncodeConfig(goConfig)
}

func writeTypeAwareDomainFixture(t *testing.T, root string) {
	t.Helper()

	testutil.WriteFile(t, root, "go.mod", "module example\n\ngo 1.24.5\n")
	testutil.WriteFile(
		t,
		root,
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
