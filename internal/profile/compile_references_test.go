package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/rulepack"
)

func TestCompileRejectsUnknownFileSetBinding(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.FileSets = nil
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown file set") {
		t.Fatalf("expected unknown file set compile error, got %v", err)
	}
}

func TestCompileRejectsUnknownLanguageBackendBinding(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Language.Backends = nil
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown language backend") {
		t.Fatalf("expected unknown language backend compile error, got %v", err)
	}
}

func TestCompileRejectsUnknownRuleBindingPathClass(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	delete(config.Paths, "go_source")
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown path class") {
		t.Fatalf("expected missing path class compile error, got %v", err)
	}
}
