package installer

import (
	"testing"

	"ciphera/tools/internal/pack/shipped"
	"ciphera/tools/internal/toolchain"
)

func TestInstallerRecognisesShippedPackInstallSpecs(t *testing.T) {
	t.Parallel()

	registry, err := shipped.DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, capability := range registry.ToolCapabilities() {
		switch capability.Install.(type) {

		case toolchain.NoInstall,
			toolchain.GoBinaryInstall,
			toolchain.NodePackageInstall,
			toolchain.ArchiveInstall:

		default:
			t.Fatalf(
				"install spec %T for tool %s is unsupported",
				capability.Install,
				capability.ID,
			)
		}
	}
}
