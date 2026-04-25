package profile

import "testing"

/* ----------------------------------------- Validation ----------------------------------------- */

func TestValidateAllowsProjectOwnedPathClasses(t *testing.T) {
	policy, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	policy.Paths.Classes["project_specific"] = []string{"internal/project/"}
	if err := policy.Validate(); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateRejectsDomainIdentifierWithoutConstructor(t *testing.T) {
	policy, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	policy.Naming.GoDomainIdentifiers["SessionKey"] = nil
	if err := policy.Validate(); err == nil {
		t.Fatal("expected empty domain identifier constructors to fail")
	}
}
