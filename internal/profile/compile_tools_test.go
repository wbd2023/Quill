package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rulepack"
)

func TestCompileRequiresActiveToolPins(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Tools = config.Tools[1:]
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "missing a tool pin") {
		t.Fatalf("expected missing tool pin error, got %v", err)
	}
}

func TestCompileRejectsUnknownToolPins(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Tools = append(config.Tools, policy.ToolPin{ID: "unknown", Version: "1.0.0"})
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("expected unknown tool pin error, got %v", err)
	}
}
