package profile

import (
	"strings"
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/rulepack"
)

func TestCompileRequiresActivePinnedTools(t *testing.T) {
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
	if err == nil || !strings.Contains(err.Error(), "missing a pinned tool") {
		t.Fatalf("expected missing pinned tool error, got %v", err)
	}
}

func TestCompileRejectsUnknownPinnedTools(t *testing.T) {
	config, err := Load(projectRoot(t))
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	config.Tools = append(config.Tools, policy.PinnedTool{ID: "unknown", Version: "1.0.0"})
	registry, err := rulepack.DefaultRegistry(config.RulePacks.Enabled)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	_, err = Compile(config, registry.Definitions())
	if err == nil || !strings.Contains(err.Error(), "unknown") {
		t.Fatalf("expected unknown pinned tool error, got %v", err)
	}
}
