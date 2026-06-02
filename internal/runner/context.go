package runner

import (
	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/toolchain"
)

type Context struct {
	RepoRoot         string
	Scope            contract.Scope
	Profile          policy.Config
	Effective        contract.EffectiveConfig
	ToolCapabilities map[string]toolchain.Capability
	ToolEnvironment  map[string]string
	GoEnvironment    map[string]string
}

func NewContext(
	repoRoot string,
	scope contract.Scope,
	config policy.Config,
	effective contract.EffectiveConfig,
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
