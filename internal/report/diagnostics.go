package report

import (
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/style"
)

// formatDiagnostic renders one finding as a compact terminal line of the form
// "file:line:column: [code] message", with location components omitted when unknown. HelpURL is
// deliberately omitted from text output: it is a machine-rendered documentation link best
// surfaced in the JSON envelope, where CI and editors can linkify it, rather than cluttering the
// human terminal line.
func formatDiagnostic(diagnostic style.Diagnostic) (line string) {
	location := diagnostic.File
	if start := diagnostic.Range.Start; start.Line > 0 {
		location = fmt.Sprintf("%s:%d", location, start.Line)
		if start.Column > 0 {
			location = fmt.Sprintf("%s:%d", location, start.Column)
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
