package shipped

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/testutil"
	"ciphera/tools/internal/testutil/profiles"
	"ciphera/tools/internal/toolchain"
)

/* ------------------------------------------- Tooling ------------------------------------------ */

func TestRegistryToolsHaveUniqueIDs(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	seenIDs := make(map[string]bool)
	for _, capability := range registry.ToolCapabilities() {
		if seenIDs[capability.ID] {
			t.Fatalf("duplicate tool ID: %s", capability.ID)
		}

		seenIDs[capability.ID] = true
	}
}

func TestRegistryToolsUseSupportedInstallStrategies(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, capability := range registry.ToolCapabilities() {
		switch install := capability.Install.(type) {

		case toolchain.NoInstall,
			toolchain.ArchiveInstall:
			_ = install

		case toolchain.GoBinaryInstall:
			if install.Source == "" {
				t.Fatalf("tool %q must define an install source", capability.ID)
			}

		case toolchain.NodePackageInstall:
			if install.Source == "" {
				t.Fatalf("tool %q must define an install source", capability.ID)
			}

		default:
			t.Fatalf(
				"tool %q uses unsupported install spec %T",
				capability.ID,
				capability.Install,
			)
		}
	}
}

func TestRegistryToolsUseSupportedVersionDetectors(t *testing.T) {
	registry, err := DefaultRegistry(nil)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	for _, capability := range registry.ToolCapabilities() {
		switch capability.Version.(type) {

		case toolchain.GoCommandVersion,
			toolchain.BuildInfoVersion,
			toolchain.PrefixedLineVersion,
			toolchain.FirstTokenVersion:

		default:
			t.Fatalf(
				"tool %q uses unsupported version spec %T",
				capability.ID,
				capability.Version,
			)
		}
	}
}

/* ---------------------------------------- Version Pins ---------------------------------------- */

func TestPinnedGoVersionMatchesModuleFiles(t *testing.T) {
	goDirectivePattern := regexp.MustCompile(`(?m)^go ([0-9]+\.[0-9]+(?:\.[0-9]+)?)$`)
	goTool := toolByID(t, tool.Go)

	rootModule := readRepoFile(t, "go.mod")
	styleModule := readRepoFile(t, filepath.Join("tools", "go.mod"))

	for _, contents := range []string{rootModule, styleModule} {
		matches := goDirectivePattern.FindStringSubmatch(contents)
		if len(matches) != 2 {
			t.Fatalf("could not find go directive in module contents:\n%s", contents)
		}

		if matches[1] != goTool.PinnedVersion {
			t.Fatalf(
				"go directive %q does not match pinned version %q",
				matches[1],
				goTool.PinnedVersion,
			)
		}
	}
}

func TestPinnedGoimportsVersionMatchesStyleModule(t *testing.T) {
	requireLinePattern := regexp.MustCompile(
		`(?m)^\s*golang\.org/x/tools (v[0-9]+\.[0-9]+\.[0-9]+)(?:$| // indirect$)`,
	)
	goimportsTool := toolByID(t, tool.Goimports)

	styleModule := readRepoFile(t, filepath.Join("tools", "go.mod"))
	matches := requireLinePattern.FindStringSubmatch(styleModule)
	if len(matches) != 2 {
		t.Fatalf("could not find golang.org/x/tools requirement in tools/go.mod")
	}

	if matches[1] != goimportsTool.PinnedVersion {
		t.Fatalf(
			"tools/go.mod pins golang.org/x/tools at %q, want %q",
			matches[1],
			goimportsTool.PinnedVersion,
		)
	}
}

/* -------------------------------------- Repository Files -------------------------------------- */

func toolByID(t *testing.T, toolID string) (tool style.Tool) {
	t.Helper()

	config := profiles.Current(t)
	registry, err := DefaultRegistry(config.EnabledPacks)
	if err != nil {
		t.Fatalf("DefaultRegistry: %v", err)
	}

	compiled, err := profile.Compile(config, registry)
	if err != nil {
		t.Fatalf("profile.Compile: %v", err)
	}

	tool, found := compiled.Effective.ToolByID(toolID)
	if !found {
		t.Fatalf("missing %s tool in registry", toolID)
	}

	return tool
}

func readRepoFile(t *testing.T, relativePath string) (contents string) {
	t.Helper()

	path := filepath.Join(testutil.RepositoryRoot(t), filepath.FromSlash(relativePath))
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", relativePath, err)
	}

	return strings.TrimSpace(string(data))
}
