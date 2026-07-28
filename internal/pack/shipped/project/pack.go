package project

import (
	"github.com/wbd2023/quill/internal/pack"
	projectpolicy "github.com/wbd2023/quill/internal/pack/shipped/project/policy"
	"github.com/wbd2023/quill/internal/pack/shipped/tool"
	"github.com/wbd2023/quill/internal/style"
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

// Pack returns the Project Shipped Pack definition. toolIDs reference the canonical Tool
// capabilities owned by the catalogue by global ID.
func Pack(toolIDs ...string) (definition pack.Definition) {
	return pack.Definition{
		ID:      PackID,
		Name:    "Project",
		ToolIDs: append([]string{}, toolIDs...),
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
