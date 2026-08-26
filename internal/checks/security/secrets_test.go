package security

import (
	"testing"

	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/testutil"
	"github.com/wbd2023/quill/internal/testutil/profiles"
)

func TestCheckSecretsFindsHighConfidenceSecretMarkers(t *testing.T) {
	root := t.TempDir()
	testutil.WriteFile(
		t,
		root,
		"internal/example/secret.txt",
		"access_key=AKI"+"AABCDEFGHIJKLMNOP\n",
	)

	result, err := CheckSecrets(root, profiles.RepositoryConfig(), style.Scope("all"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Diagnostics) == 0 {
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
	root := t.TempDir()
	testutil.WriteFile(t, root, "internal/example/doc.txt", "ordinary content\n")

	result, err := CheckSecrets(root, profiles.RepositoryConfig(), style.Scope("all"))
	if err != nil {
		t.Fatalf("expected committed-secret check to pass, diagnostics: %#v", result.Diagnostics)
	}
}
