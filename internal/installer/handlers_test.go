package installer

import (
	"io"
	"strings"
	"testing"

	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func TestInstallToolRejectsUnknownInstallKind(t *testing.T) {
	err := installTool(
		runtime.LayoutForToolsDir(t.TempDir()),
		io.Discard,
		style.Tool{ID: "example"},
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

func TestInstallerSupportsShippedPackToolInstallKinds(t *testing.T) {
	registry, err := shipped.DefaultRegistry(nil)
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
	}
}
