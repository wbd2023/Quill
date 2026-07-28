package pack

import "github.com/wbd2023/quill/internal/toolchain"

// CloneCapability returns a defensive copy of capability, including the mutable installer data a
// Capability may carry (for example the platform map of a GitHub release install). Callers must
// never observe aliasing into a catalogue's canonical capabilities.
func CloneCapability(capability toolchain.Capability) (clone toolchain.Capability) {
	clone = capability
	clone.Install = CloneInstallMethod(capability.Install)
	return clone
}

// CloneCapabilities returns defensive deep copies of capabilities.
func CloneCapabilities(capabilities []toolchain.Capability) (clones []toolchain.Capability) {
	clones = make([]toolchain.Capability, len(capabilities))
	for index, capability := range capabilities {
		clones[index] = CloneCapability(capability)
	}
	return clones
}

// CloneInstallMethod returns a defensive copy of method. Only install methods that carry mutable
// state (a platform map) are deep copied; the remainder are value types returned unchanged.
func CloneInstallMethod(method toolchain.InstallMethod) (clone toolchain.InstallMethod) {
	switch install := method.(type) {
	case toolchain.GitHubInstall:
		install.Platforms = clonePlatformMap(install.Platforms)
		return install

	case *toolchain.GitHubInstall:
		// GitHubInstall has a value receiver, so a pointer also satisfies InstallMethod. Clone
		// the pointed-to value and its platform map so callers cannot mutate the original through
		// either form.
		if install == nil {
			return method
		}

		duplicated := *install
		duplicated.Platforms = clonePlatformMap(install.Platforms)
		return &duplicated

	default:
		return method
	}
}

func clonePlatformMap(platforms map[string]string) (cloned map[string]string) {
	if platforms == nil {
		return nil
	}

	cloned = make(map[string]string, len(platforms))
	for platform, token := range platforms {
		cloned[platform] = token
	}
	return cloned
}
