package builtin

import (
	"ciphera/tools/internal/pack/builtin/bash"
	"ciphera/tools/internal/pack/builtin/golang"
	"ciphera/tools/internal/pack/builtin/markdown"
	"ciphera/tools/internal/pack/builtin/project"
	"ciphera/tools/internal/pack/builtin/security"
	"ciphera/tools/internal/pack/builtin/text"
	"ciphera/tools/internal/pack/builtin/vocabulary"
)

const (
	PackProject    = project.PackID
	PackText       = text.PackID
	PackMarkdown   = markdown.PackID
	PackBash       = bash.PackID
	PackGo         = golang.PackID
	PackSecurity   = security.PackID
	PackVocabulary = vocabulary.PackID
)
