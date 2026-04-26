package rulepack

func markdownPack() (pack Pack) {
	return Pack{
		ID:   PackMarkdown,
		Name: "Markdown",
		Tools: selectTools(
			ToolMarkdownlint,
		),
		Rules: []RuleDefinition{
			fileCommandRuleWithConfig(
				"markdown/style",
				"Markdown style",
				ToolMarkdownlint,
				"markdown",
				nil,
				"-c",
				".markdownlint.jsonc",
			),
		},
	}
}
