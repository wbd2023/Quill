package runner

import (
	"path/filepath"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/profile"
	"ciphera/tools/internal/runtime"
)

type Context struct {
	RepoRoot        string
	Scope           contract.Scope
	Policy          profile.Profile
	Effective       profile.EffectiveConfig
	Layout          runtime.Layout
	ToolEnvironment map[string]string
	GoEnvironment   map[string]string
}

func NewContext(
	repoRoot string,
	scope contract.Scope,
	policy profile.Profile,
	effective profile.EffectiveConfig,
) (context Context) {
	layout := runtime.LayoutForRepository(repoRoot)
	goEnvironment := layout.GoEnvironment()
	goEnvironment["GOLANGCI_LINT_CACHE"] = filepath.Join(layout.CacheDir, "golangci")

	return Context{
		RepoRoot:        repoRoot,
		Scope:           scope,
		Policy:          policy,
		Effective:       effective,
		Layout:          layout,
		ToolEnvironment: layout.ToolEnvironment(),
		GoEnvironment:   goEnvironment,
	}
}
