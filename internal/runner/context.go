package runner

import (
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

// Context is context.
type Context struct {
	RepoRoot         string
	Scope            style.Scope
	Profile          policy.Config
	Effective        style.EffectiveConfig
	ToolCapabilities map[string]toolchain.Capability
	ToolEnvironment  map[string]string
	GoEnvironment    map[string]string
}

// NewContext new context.
func NewContext(
	repoRoot string,
	scope style.Scope,
	config policy.Config,
	effective style.EffectiveConfig,
	capabilities []toolchain.Capability,
	toolEnvironment map[string]string,
	goEnvironment map[string]string,
) (context Context) {
	return Context{
		RepoRoot:         repoRoot,
		Scope:            scope,
		Profile:          config,
		Effective:        effective,
		ToolCapabilities: toolchain.CapabilitiesByID(capabilities),
		ToolEnvironment:  toolEnvironment,
		GoEnvironment:    goEnvironment,
	}
}
