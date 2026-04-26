package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
)

/* ----------------------------------------- Validation ----------------------------------------- */

func TestValidateAllowsProjectOwnedPathClasses(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Paths.Classes["project_specific"] = []string{"internal/project/"}
	if err := Validate(config); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateRejectsDomainIdentifierWithoutConstructor(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Naming.GoDomainIdentifiers["SessionKey"] = nil
	if err := Validate(config); err == nil {
		t.Fatal("expected empty domain identifier constructors to fail")
	}
}

func TestValidateRequiresCurrentSchemaVersion(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.SchemaVersion = 2
	if err := Validate(config); err == nil || !strings.Contains(err.Error(), "version 2") {
		t.Fatalf("expected schema version error, got %v", err)
	}
}

func TestValidateRejectsUnknownDefaultScope(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Repository.DefaultScope = "unknown"
	if err := Validate(config); err == nil || !strings.Contains(err.Error(), "default_scope") {
		t.Fatalf("expected unknown default scope error, got %v", err)
	}
}

func TestValidateRejectsInvalidSectionHeaderPolicy(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Formatting.SectionHeaders.ShortFileMaxLines = 100
	err = Validate(config)
	if err == nil || !strings.Contains(err.Error(), "short_file_max_lines") {
		t.Fatalf("expected invalid section header policy error, got %v", err)
	}
}

func TestValidateRejectsUnknownFileSetScope(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.FileSets[0].Prefixes[contract.Scope("unknown")] = []string{"unknown/"}
	if err := Validate(config); err == nil || !strings.Contains(err.Error(), "unknown scope") {
		t.Fatalf("expected unknown file-set scope error, got %v", err)
	}
}

func TestValidateRejectsUnknownRuleScope(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Rules[0].Scope = "unknown"
	if err := Validate(config); err == nil || !strings.Contains(err.Error(), "unknown scope") {
		t.Fatalf("expected unknown rule scope error, got %v", err)
	}
}
