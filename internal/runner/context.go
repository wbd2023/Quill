package runner

import (
	"ciphera/tools/internal/lockfile"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

// Context carries loaded profile and toolchain state through a check or install run.
type Context struct {
	RepoRoot         string
	Scope            style.Scope
	Profile          policy.Config
	Effective        style.EffectiveConfig
	ToolCapabilities map[string]toolchain.Capability
	ToolEnvironment  map[string]string
	GoEnvironment    map[string]string
	Lockfile         lockfile.Lockfile
}

// NewContext constructs a Context from loaded profile and toolchain state.
func NewContext(
	repoRoot string,
	scope style.Scope,
	config policy.Config,
	effective style.EffectiveConfig,
	capabilities []toolchain.Capability,
	toolEnvironment map[string]string,
	goEnvironment map[string]string,
	lockfile lockfile.Lockfile,
) (context Context) {
	return Context{
		RepoRoot:         repoRoot,
		Scope:            scope,
		Profile:          config,
		Effective:        effective,
		ToolCapabilities: toolchain.CapabilitiesByID(capabilities),
		ToolEnvironment:  toolEnvironment,
		GoEnvironment:    goEnvironment,
		Lockfile:         lockfile,
	}
}
