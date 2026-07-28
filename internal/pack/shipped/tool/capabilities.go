package tool

import (
	"github.com/wbd2023/quill/internal/toolchain"
)

// BuildAll returns the canonical Shipped Tool capabilities. Each Tool is defined exactly once here;
// Packs reference these Tools by global ID (see ids.go) and the catalogue resolves the references,
// so a Tool capability is never duplicated across Packs.
func BuildAll() (capabilities []toolchain.Capability) {
	return []toolchain.Capability{
		buildBuiltin(Go, "Go", "go",
			toolchain.DetectByCommand("version", toolchain.ExtractGoToken)),
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
