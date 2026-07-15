package shipped

import (
	"testing"

	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/markdown"
	"ciphera/tools/internal/pack/shipped/project"
	"ciphera/tools/internal/pack/shipped/text"
	"ciphera/tools/internal/pack/shipped/tool"
)

func TestDefaultRegistryLoadsEnabledPacks(t *testing.T) {
	registry, err := DefaultRegistry([]string{project.PackID, text.PackID, markdown.PackID})
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	if len(registry.Packs()) != 3 {
		t.Fatalf("expected 3 packs, got %d", len(registry.Packs()))
	}

	if len(registry.Rules()) == 0 {
		t.Fatal("expected enabled packs to register rules")
	}

	capabilities := registry.ToolCapabilities()
	found := false
	for _, capability := range capabilities {
		if capability.ID == tool.Markdownlint {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected markdown pack tool to be registered")
	}
}

func TestDefaultRegistryRejectsUnknownPack(t *testing.T) {
	if _, err := DefaultRegistry([]string{"unknown"}); err == nil {
		t.Fatal("expected unknown pack to be rejected")
	}
}

func TestDefaultRegistryRejectsDuplicatePack(t *testing.T) {
	if _, err := DefaultRegistry([]string{golang.PackID, golang.PackID}); err == nil {
		t.Fatal("expected duplicate pack to be rejected")
	}
}
