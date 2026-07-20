package execution

import (
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

// RunContext carries loaded profile and toolchain state through a check or install run.
type RunContext struct {
	RepoRoot        string
	Scope           style.Scope
	Profile         policy.Config
	Effective       style.Plan
	Tools           map[string]toolchain.Tool
	ToolEnvironment map[string]string
	GoEnvironment   map[string]string
}

// NewRunContext constructs a RunContext from loaded profile and toolchain state. It joins each
// capability with its pinned version and execution limits into a toolchain.Tool.
func NewRunContext(
	repoRoot string,
	scope style.Scope,
	config policy.Config,
	effective style.Plan,
	capabilities []toolchain.Capability,
	toolEnvironment map[string]string,
	goEnvironment map[string]string,
) (context RunContext) {
	tools := make(map[string]toolchain.Tool, len(capabilities))
	for _, capability := range capabilities {
		pin, _ := config.Tools.Lookup(capability.ID)
		tools[capability.ID] = toolchain.Tool{
			ID:               capability.ID,
			Name:             capability.Name,
			PinnedVersion:    pin.Version,
			TimeoutSeconds:   pin.TimeoutSeconds,
			OutputLimitBytes: pin.OutputLimitBytes,
			Command:          capability.Command,
			Version:          capability.Version,
			Install:          capability.Install,
		}
	}

	return RunContext{
		RepoRoot:        repoRoot,
		Scope:           scope,
		Profile:         config,
		Effective:       effective,
		Tools:           tools,
		ToolEnvironment: toolEnvironment,
		GoEnvironment:   goEnvironment,
	}
}
