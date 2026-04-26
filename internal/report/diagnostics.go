package report

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/contract"
)

func formatDiagnostic(diagnostic contract.Diagnostic) (line string) {
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

func groupLabel(group contract.RuleGroup) (label string) {
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
