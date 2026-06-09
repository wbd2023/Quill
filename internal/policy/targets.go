package policy

import "ciphera/tools/internal/style"

// TargetConfigs defines the targets available to rule bindings.
type TargetConfigs []TargetConfig

// TargetConfig binds language-specific target settings to a repository scope.
type TargetConfig struct {
	Name             string
	Language         string
	Scope            style.Scope
	WorkingDirectory string
	FormatPaths      []string
	CheckPaths       []string
}

// Lookup returns the named target.
func (targets TargetConfigs) Lookup(name string) (target TargetConfig, found bool) {
	for _, candidate := range targets {
		if candidate.Name == name {
			return candidate, true
		}
	}

	return TargetConfig{}, false
}
