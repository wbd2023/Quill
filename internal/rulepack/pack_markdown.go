package rulepack

import "ciphera/tools/internal/contract"

/* ---------------------------------------- Markdown Pack --------------------------------------- */

func markdownPack() (pack Pack) {
	return Pack{
		ID:   PackMarkdown,
		Name: "Markdown",
		Tools: selectTools(
			contract.ToolMarkdownlint,
		),
		Rules: []RuleDefinition{
			fileCommandRuleWithConfig(
				"markdown/style",
				"Markdown style",
				contract.ToolMarkdownlint,
				"markdown",
				nil,
				"-c",
				".markdownlint.jsonc",
			),
		},
	}
}
