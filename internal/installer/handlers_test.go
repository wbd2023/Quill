package installer

import (
	"testing"

	"github.com/wbd2023/Quill/internal/pack/shipped"
	"github.com/wbd2023/Quill/internal/toolchain"
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
			toolchain.GoInstall,
			toolchain.NpmInstall,
			toolchain.GitHubInstall:

		default:
			t.Fatalf(
				"install spec %T for tool %s is unsupported",
				capability.Install,
				capability.ID,
			)
		}
	}
}
