package project

import (
	"ciphera/tools/internal/checks/projectpolicy"
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

// PackID is the canonical identifier for this Pack.
const PackID = "project"

// pack constants.
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
			Validate: projectpolicy.ValidatePackConfig,
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
			tool.Go,
			tool.Goimports,
			tool.Misspell,
			tool.GolangciLint,
			tool.Shfmt,
			tool.Shellcheck,
			tool.Markdownlint,
		),
		projectRule(
			"profile/enforcement-levels",
			"Enforcement levels",
			CheckEnforcementLevels,
		),
		projectRule(
			"profile/quality-commands",
			"Quality commands",
			CheckCommands,
		),
		projectRule(
			"profile/excluded-directories",
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
		Check: style.ToolchainExecution{
			ToolIDs: append([]string{}, toolIDs...),
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
		Check: style.ProfileExecution{
			Check: check,
		},
	}
}
