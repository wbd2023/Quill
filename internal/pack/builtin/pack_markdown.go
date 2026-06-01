package builtin

import "ciphera/tools/internal/policy"

func markdownPack() (pack Pack) {
	return Pack{
		ID:   PackMarkdown,
		Name: "Markdown",
		Tools: selectTools(
			ToolMarkdownlint,
		),
		FileSets: markdownFileSets(),
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

func markdownFileSets() (fileSets policy.FileSets) {
	return append(fileSets, policy.FileSetConfig{
		Name: "markdown",
		Include: policy.FileSetInclude{
			Extensions: []string{".md"},
		},
	})
}
