package rulepack

import (
	"slices"

	"ciphera/tools/internal/toolchain"
)

func coreTools() (tools []toolchain.Capability) {
	return []toolchain.Capability{
		builtinTool(ToolGo, "Go", "go", ToolVersionGoCommand),
		goBinaryTool(
			ToolGoimports,
			"goimports",
			"goimports",
			"golang.org/x/tools",
			"golang.org/x/tools/cmd/goimports",
		),
		goBinaryTool(
			ToolMisspell,
			"misspell",
			"misspell",
			"github.com/client9/misspell",
			"github.com/client9/misspell/cmd/misspell",
		),
		goBinaryTool(
			ToolGolangciLint,
			"golangci-lint",
			"golangci-lint",
			"github.com/golangci/golangci-lint/v2",
			"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
		),
		goBinaryTool(
			ToolShfmt,
			"shfmt",
			"shfmt",
			"mvdan.cc/sh/v3",
			"mvdan.cc/sh/v3/cmd/shfmt",
		),
		shellcheckArchiveTool(),
		nodePackageTool(
			ToolMarkdownlint,
			"markdownlint",
			"markdownlint",
			"markdownlint-cli",
		),
	}
}

func selectTools(toolIDs ...string) (tools []toolchain.Capability) {
	wanted := make(map[string]bool, len(toolIDs))
	for _, toolID := range toolIDs {
		wanted[toolID] = true
	}

	for _, tool := range coreTools() {
		if wanted[tool.ID] {
			tools = append(tools, tool)
		}
	}

	slices.SortFunc(tools, func(left toolchain.Capability, right toolchain.Capability) int {
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
