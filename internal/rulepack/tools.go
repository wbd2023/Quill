package rulepack

import (
	"slices"

	"ciphera/tools/internal/contract"
)

/* -------------------------------------------- Tools ------------------------------------------- */

func coreTools() (tools []contract.Tool) {
	return []contract.Tool{
		builtinTool(contract.ToolGo, "Go", "go", "1.24.5", contract.ToolVersionGoCommand),
		goBinaryTool(
			contract.ToolGoimports,
			"goimports",
			"goimports",
			"v0.42.0",
			"golang.org/x/tools",
			"golang.org/x/tools/cmd/goimports",
		),
		goBinaryTool(
			contract.ToolMisspell,
			"misspell",
			"misspell",
			"v0.3.4",
			"github.com/client9/misspell",
			"github.com/client9/misspell/cmd/misspell",
		),
		goBinaryTool(
			contract.ToolGolangciLint,
			"golangci-lint",
			"golangci-lint",
			"v2.6.2",
			"github.com/golangci/golangci-lint/v2",
			"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
		),
		goBinaryTool(
			contract.ToolShfmt,
			"shfmt",
			"shfmt",
			"v3.12.0",
			"mvdan.cc/sh/v3",
			"mvdan.cc/sh/v3/cmd/shfmt",
		),
		shellcheckArchiveTool(),
		nodePackageTool(
			contract.ToolMarkdownlint,
			"markdownlint",
			"markdownlint",
			"0.45.0",
			"markdownlint-cli",
		),
	}
}

func selectTools(toolIDs ...string) (tools []contract.Tool) {
	wanted := make(map[string]bool, len(toolIDs))
	for _, toolID := range toolIDs {
		wanted[toolID] = true
	}

	for _, tool := range coreTools() {
		if wanted[tool.ID] {
			tools = append(tools, tool)
		}
	}

	slices.SortFunc(tools, func(left contract.Tool, right contract.Tool) int {
		if left.ID < right.ID {
			return -1
		}

		if left.ID > right.ID {
			return 1
		}

		return 0
	})
	return tools
}
