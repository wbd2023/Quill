package security

import (
	"testing"

	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
)

func TestCheckSecretsFindsHighConfidenceSecretMarkers(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(
		t,
		repoRoot,
		"internal/example/secret.txt",
		"access_key=AKI"+"AABCDEFGHIJKLMNOP\n",
	)

	result, err := CheckSecrets(repoRoot, profiles.RepositoryConfig(t), style.Scope("all"))
	if err == nil {
		t.Fatal("expected committed-secret failure")
	}

	if !hasDiagnostic(
		result,
		"security/secrets/aws-key",
		"internal/example/secret.txt",
		1,
		"possible AWS access key",
	) {
		t.Fatalf("expected token diagnostic, got: %#v", result.Diagnostics)
	}
}

func TestCheckSecretsPassesOrdinaryFiles(t *testing.T) {
	repoRoot := t.TempDir()
	testutil.WriteFile(t, repoRoot, "internal/example/doc.txt", "ordinary content\n")

	result, err := CheckSecrets(repoRoot, profiles.RepositoryConfig(t), style.Scope("all"))
	if err != nil {
		t.Fatalf("expected committed-secret check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
