package repostyle

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/fixtures"
	"ciphera/tools/internal/fixtures/profiles"
)

func TestCheckSecretsFindsHighConfidenceSecretMarkers(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(
		t,
		repoRoot,
		"internal/example/secret.txt",
		"access_key=AKI"+"AABCDEFGHIJKLMNOP\n",
	)

	output, err := CheckSecrets(repoRoot, profiles.RepositoryConfig(t), contract.ScopeAll)
	if err == nil {
		t.Fatal("expected committed-secret failure")
	}

	if !strings.Contains(output, "possible AWS access key") {
		t.Fatalf("expected token violation, got:\n%s", output)
	}
}

func TestCheckSecretsPassesOrdinaryFiles(t *testing.T) {
	repoRoot := t.TempDir()
	fixtures.WriteFile(t, repoRoot, "internal/example/doc.txt", "ordinary content\n")

	output, err := CheckSecrets(repoRoot, profiles.RepositoryConfig(t), contract.ScopeAll)
	if err != nil {
		t.Fatalf("expected committed-secret check to pass, output:\n%s", output)
	}
}
