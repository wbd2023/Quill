package pack

import (
	"testing"

	"github.com/wbd2023/quill/internal/toolchain"
)

func TestCloneCapabilitiesReturnsIndependentCopies(t *testing.T) {
	original := []toolchain.Capability{
		{
			ID:      "shellcheck",
			Name:    "shellcheck",
			Command: "shellcheck",
			Install: toolchain.GitHubInstall{
				Owner:      "koalaman",
				Repository: "shellcheck",
				Platforms: map[string]string{
					"linux/amd64": "linux.x86_64",
				},
			},
		},
	}

	clones := CloneCapabilities(original)
	clones[0].Install.(toolchain.GitHubInstall).Platforms["linux/amd64"] = "mutated"

	platforms := original[0].Install.(toolchain.GitHubInstall).Platforms
	if got := platforms["linux/amd64"]; got != "linux.x86_64" {
		t.Fatalf("original platform map mutated via clone: %q", got)
	}
}

func TestCloneInstallMethodPreservesNilPlatformMap(t *testing.T) {
	install := toolchain.GitHubInstall{Owner: "o", Repository: "r"}
	cloned := CloneInstallMethod(install).(toolchain.GitHubInstall)

	if cloned.Platforms != nil {
		t.Fatalf("expected nil platform map preserved, got %#v", cloned.Platforms)
	}
}

func TestCloneInstallMethodReturnsValueInstallers(t *testing.T) {
	cases := []toolchain.InstallMethod{
		toolchain.NoInstall{},
		toolchain.GoInstall{Source: "example.com/mod"},
		toolchain.NpmInstall{Source: "pkg"},
	}

	for _, method := range cases {
		if CloneInstallMethod(method) == nil {
			t.Fatalf("CloneInstallMethod returned nil for %#v", method)
		}
	}
}

func TestCloneInstallMethodDefendsPointerGitHubInstall(t *testing.T) {
	original := &toolchain.GitHubInstall{
		Owner:      "koalaman",
		Repository: "shellcheck",
		Platforms:  map[string]string{"linux/amd64": "linux.x86_64"},
	}

	cloned := CloneInstallMethod(original).(*toolchain.GitHubInstall)

	// Mutating the clone's platform map and owner must not touch the original.
	cloned.Platforms["linux/amd64"] = "mutated"
	cloned.Owner = "changed"

	if got := original.Platforms["linux/amd64"]; got != "linux.x86_64" {
		t.Fatalf("original platform map mutated via pointer clone: %q", got)
	}
	if got := original.Owner; got != "koalaman" {
		t.Fatalf("original owner mutated via pointer clone: %q", got)
	}
}

func TestCloneCapabilitiesDefendsPointerInstaller(t *testing.T) {
	capability := toolchain.Capability{
		ID:      "shellcheck",
		Name:    "shellcheck",
		Command: "shellcheck",
		Install: &toolchain.GitHubInstall{Platforms: map[string]string{"linux/amd64": "x"}},
	}

	clones := CloneCapabilities([]toolchain.Capability{capability})
	clones[0].Install.(*toolchain.GitHubInstall).Platforms["linux/amd64"] = "mutated"

	if got := capability.Install.(*toolchain.GitHubInstall).Platforms["linux/amd64"]; got != "x" {
		t.Fatalf("original capability mutated via cloned pointer installer: %q", got)
	}
}
