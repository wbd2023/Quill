package project

import (
	projectrules "ciphera/tools/internal/checks/project"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

const (
	PackID = "project"

	ToolGo           = "go"
	ToolGoimports    = "goimports"
	ToolMisspell     = "misspell"
	ToolGolangciLint = "golangci-lint"
	ToolShfmt        = "shfmt"
	ToolShellcheck   = "shellcheck"
	ToolMarkdownlint = "markdownlint"
)

const (
	CheckEnforcementLevels   = "enforcement_levels"
	CheckExcludedDirectories = "excluded_directories"
	CheckCommands            = "commands"
)

const ruleGroupProject style.RuleGroup = "project"

// Pack returns the Project Shipped Pack definition.
func Pack(tools []toolchain.Capability) (definition pack.Definition) {
	return pack.Definition{
		ID:    PackID,
		Name:  "Project",
		Tools: append([]toolchain.Capability{}, tools...),
		Config: pack.Config{
			Required: true,
			Validate: projectrules.ValidatePackConfig,
		},
		Rules: rules(),
	}
}

/* ----------------------------------------- Rule Lists ----------------------------------------- */

func rules() (rules []style.RuleDefinition) {
	return []style.RuleDefinition{
		toolchainRule(
			"toolchain/check-versions",
			"Pinned toolchain versions",
			ToolGo,
			ToolGoimports,
			ToolMisspell,
			ToolGolangciLint,
			ToolShfmt,
			ToolShellcheck,
			ToolMarkdownlint,
		),
		projectRule(
			"project/enforcement-levels",
			"Enforcement levels",
			CheckEnforcementLevels,
		),
		projectRule(
			"project/quality-commands",
			"Quality commands",
			CheckCommands,
		),
		projectRule(
			"project/excluded-directories",
			"Excluded directories",
			CheckExcludedDirectories,
		),
	}
}

/* ---------------------------------------- Rule Builders --------------------------------------- */

func toolchainRule(
	id string,
	name string,
	toolIDs ...string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupProject,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionToolchain,
			Detail: style.ToolchainExecution{
				ToolIDs: append([]string{}, toolIDs...),
			},
		},
	}
}

func projectRule(
	id string,
	name string,
	check string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupProject,
		Check: style.ExecutionSpec{
			Kind: style.ExecutionProject,
			Detail: style.ProjectExecution{
				Check: check,
			},
		},
	}
}
