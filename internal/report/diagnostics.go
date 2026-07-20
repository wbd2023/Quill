package report

import (
	"fmt"
	"strings"

	"github.com/wbd2023/Quill/internal/style"
)

func formatDiagnostic(diagnostic style.Diagnostic) (line string) {
	location := diagnostic.File
	if diagnostic.Line > 0 {
		location = fmt.Sprintf("%s:%d", location, diagnostic.Line)
		if diagnostic.Column > 0 {
			location = fmt.Sprintf("%s:%d", location, diagnostic.Column)
		}
	}

	if diagnostic.Code == "" {
		return fmt.Sprintf("%s %s", location, diagnostic.Message)
	}

	return fmt.Sprintf("%s: [%s] %s", location, diagnostic.Code, diagnostic.Message)
}

func groupLabel(group style.RuleGroup) (label string) {
	words := strings.FieldsFunc(string(group), func(rune rune) bool {
		return rune == '_' || rune == '-' || rune == '/'
	})
	for index, word := range words {
		if word == "" {
			continue
		}

		words[index] = strings.ToUpper(word[:1]) + word[1:]
	}

	return strings.Join(words, " ")
}
