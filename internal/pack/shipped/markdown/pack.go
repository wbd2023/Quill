package markdown

import (
	"github.com/wbd2023/Quill/internal/pack"
	"github.com/wbd2023/Quill/internal/pack/shipped/tool"
	"github.com/wbd2023/Quill/internal/policy"
	"github.com/wbd2023/Quill/internal/style"
	"github.com/wbd2023/Quill/internal/toolchain"
)

// PackID is the canonical identifier for this Pack.
const PackID = "markdown"

const ruleGroupExternal style.RuleGroup = "external_tools"

// Pack returns the Markdown Shipped Pack definition.
func Pack(tools []toolchain.Capability) (definition pack.Definition) {
	return pack.Definition{
		ID:       PackID,
		Name:     "Markdown",
		Tools:    append([]toolchain.Capability{}, tools...),
		FileSets: fileSets(),
		Rules: []style.RuleDefinition{
			fileCommandRuleWithConfig(
				"markdown/style",
				"Markdown style",
				tool.Markdownlint,
				"markdown",
				nil,
				"-c",
				".markdownlint.jsonc",
			),
		},
	}
}

func fileSets() (fileSets policy.FileSets) {
	return append(fileSets, policy.FileSetConfig{
		Name: "markdown",
		Include: policy.FileSetInclude{
			Extensions: []string{".md"},
		},
	})
}

func fileCommandRuleWithConfig(
	id string,
	name string,
	toolID string,
	fileSet string,
	arguments []string,
	configArgument string,
	configFile string,
) (rule style.RuleDefinition) {
	rule = fileCommandRule(id, name, toolID, fileSet, arguments)
	execution := rule.Check.(style.FileCommandExecution)
	execution.ConfigArgument = configArgument
	execution.ConfigFile = configFile
	rule.Check = execution
	return rule
}

func fileCommandRule(
	id string,
	name string,
	toolID string,
	fileSet string,
	arguments []string,
) (rule style.RuleDefinition) {
	return style.RuleDefinition{
		ID:    id,
		Name:  name,
		Group: ruleGroupExternal,
		Check: style.FileCommandExecution{
			ToolID:    toolID,
			FileSet:   fileSet,
			Arguments: append([]string{}, arguments...),
		},
	}
}
