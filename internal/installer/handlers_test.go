package installer

import (
	"strings"
	"testing"

	"io"

	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/runtime"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

func TestInstallToolRejectsUnknownInstallKind(t *testing.T) {
	err := installTool(
		runtime.NewLayout(t.TempDir()),
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

func TestInstallerRecognisesShippedPackToolInstallKinds(t *testing.T) {
	registry, err := shipped.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, capability := range registry.ToolCapabilities() {
		switch capability.InstallKind {

		case toolchain.InstallKindNone,
			toolchain.InstallKindGoBinary,
			toolchain.InstallKindNodePackage,
			toolchain.InstallKindArchive:

		default:
			t.Fatalf(
				"install kind %q for tool %s is unsupported",
				capability.InstallKind,
				capability.ID,
			)
		}
	}
}
