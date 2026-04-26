package runtime

import (
	"io"
	"strings"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/rulepack"
	"ciphera/tools/internal/toolchain"
)

func TestInstallToolRejectsUnknownInstallKind(t *testing.T) {
	err := installTool(
		LayoutForToolsDir(t.TempDir()),
		io.Discard,
		contract.Tool{ID: "example"},
		toolchain.Capability{
			ID:          "example",
			InstallKind: "unknown",
		},
	)
	if err == nil {
		t.Fatal("expected unknown install kind to fail")
	}

	if !strings.Contains(err.Error(), "unsupported install strategy") {
		t.Fatalf("unexpected install error: %v", err)
	}
}

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

func TestRuntimeSupportsRulepackToolCapabilityKinds(t *testing.T) {
	registry, err := rulepack.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, capability := range registry.ToolCapabilities() {
		if !SupportsInstallKind(capability.InstallKind) {
			t.Fatalf(
				"install kind %q for tool %s is unsupported",
				capability.InstallKind,
				capability.ID,
			)
		}

		if !SupportsVersionKind(capability.VersionKind) {
			t.Fatalf(
				"version kind %q for tool %s is unsupported",
				capability.VersionKind,
				capability.ID,
			)
		}
	}
}
