package tool

import (
	"slices"

	"ciphera/tools/internal/toolchain"
)

// BuildAll build all.
func BuildAll() (capabilities []toolchain.Capability) {
	return []toolchain.Capability{
		buildBuiltin(Go, "Go", "go", toolchain.GoCommandVersion{}),
		buildGoBinary(
			Goimports,
			"goimports",
			"goimports",
			"golang.org/x/tools",
			"golang.org/x/tools/cmd/goimports",
		),
		buildGoBinary(
			Misspell,
			"misspell",
			"misspell",
			"github.com/client9/misspell",
			"github.com/client9/misspell/cmd/misspell",
		),
		buildGoBinary(
			GolangciLint,
			"golangci-lint",
			"golangci-lint",
			"github.com/golangci/golangci-lint/v2",
			"github.com/golangci/golangci-lint/v2/cmd/golangci-lint",
		),
		buildGoBinary(
			Shfmt,
			"shfmt",
			"shfmt",
			"mvdan.cc/sh/v3",
			"mvdan.cc/sh/v3/cmd/shfmt",
		),
		buildShellcheckArchive(),
		buildNodePackage(
			Markdownlint,
			"markdownlint",
			"markdownlint",
			"markdownlint-cli",
		),
	}
}

// Select returns only the capabilities whose IDs match the given tool IDs.
func Select(toolIDs ...string) (capabilities []toolchain.Capability) {
	wanted := make(map[string]bool, len(toolIDs))
	for _, toolID := range toolIDs {
		wanted[toolID] = true
	}

	for _, capability := range BuildAll() {
		if wanted[capability.ID] {
			capabilities = append(capabilities, capability)
		}
	}

	slices.SortFunc(capabilities, func(left toolchain.Capability, right toolchain.Capability) int {
		if left.ID < right.ID {
			return -1
		}

		if left.ID > right.ID {
			return 1
		}

		return 0
	})
	return capabilities
}
