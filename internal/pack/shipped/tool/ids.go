package tool

import (
	"ciphera/tools/internal/pack/shipped/bash"
	"ciphera/tools/internal/pack/shipped/golang"
	"ciphera/tools/internal/pack/shipped/markdown"
	"ciphera/tools/internal/pack/shipped/text"
)

const (
	Go           = golang.ToolGo
	Goimports    = golang.ToolGoimports
	Misspell     = text.ToolMisspell
	GolangciLint = golang.ToolGolangciLint
	Shfmt        = bash.ToolShfmt
	Shellcheck   = bash.ToolShellcheck
	Markdownlint = markdown.ToolMarkdownlint
)
