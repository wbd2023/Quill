package markdown

import (
	"ciphera/tools/internal/pack"
	"ciphera/tools/internal/pack/shipped/tool"
	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/style"
	"ciphera/tools/internal/toolchain"
)

// PackID is pack i d.
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
	execution := rule.Check.Detail.(style.FileCommandExecution)
	execution.ConfigArgument = configArgument
	execution.ConfigFile = configFile
	rule.Check.Detail = execution
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
		Check: style.ExecutionSpec{
			Kind: style.ExecutionFileCommand,
			Detail: style.FileCommandExecution{
				ToolID:    toolID,
				FileSet:   fileSet,
				Arguments: append([]string{}, arguments...),
			},
		},
	}
}
