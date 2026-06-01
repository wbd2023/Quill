package builtin

import "testing"

func TestDefaultRegistryLoadsEnabledPacks(t *testing.T) {
	registry, err := DefaultRegistry([]string{PackProject, PackText, PackMarkdown})
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	if len(registry.Packs()) != 3 {
		t.Fatalf("expected 3 packs, got %d", len(registry.Packs()))
	}

	if len(registry.Rules()) == 0 {
		t.Fatal("expected enabled packs to register rules")
	}

	if _, found := registry.ToolByID(ToolMarkdownlint); !found {
		t.Fatal("expected markdown pack tool to be registered")
	}
}

func TestDefaultRegistryRejectsUnknownPack(t *testing.T) {
	if _, err := DefaultRegistry([]string{"unknown"}); err == nil {
		t.Fatal("expected unknown pack to be rejected")
	}
}

func TestDefaultRegistryRejectsDuplicatePack(t *testing.T) {
	if _, err := DefaultRegistry([]string{PackGo, PackGo}); err == nil {
		t.Fatal("expected duplicate pack to be rejected")
	}
}
