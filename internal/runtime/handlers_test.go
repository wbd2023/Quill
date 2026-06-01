package runtime

import (
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/pack/builtin"
	"ciphera/tools/internal/toolchain"
)

func TestDetectVersionRejectsUnknownVersionKind(t *testing.T) {
	_, err := detectVersion(
		contract.Tool{ID: "example"},
		toolchain.Capability{
			ID:          "example",
			VersionKind: "unknown",
		},
		"/bin/true",
		nil,
	)
	if err == nil {
		t.Fatal("expected unknown version kind to fail")
	}

	if !strings.Contains(err.Error(), "unsupported version detector") {
		t.Fatalf("unexpected version error: %v", err)
	}
}

func TestRuntimeSupportsRulepackToolVersionKinds(t *testing.T) {
	registry, err := builtin.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, capability := range registry.ToolCapabilities() {
		if !SupportsVersionKind(capability.VersionKind) {
			t.Fatalf(
				"version kind %q for tool %s is unsupported",
				capability.VersionKind,
				capability.ID,
			)
		}
	}
}
