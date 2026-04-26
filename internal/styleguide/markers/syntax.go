package markers

import "strings"

const (
	markerPrefix    = "style: "
	reasonSeparator = " because: "
)

const (
	shellCommentPrefix     = "# "
	goLineCommentPrefix    = "// "
	blockCommentPrefix     = "/* "
	blockLineCommentPrefix = "* "
	blockCommentSuffix     = "*/"
)

type commentForm struct {
	prefix string
	suffix string
}

func extractDirective(line string) (directive string, found bool) {
	forms := [...]commentForm{
		{prefix: shellCommentPrefix},
		{prefix: goLineCommentPrefix},
		{prefix: blockCommentPrefix, suffix: blockCommentSuffix},
		{prefix: blockLineCommentPrefix, suffix: blockCommentSuffix},
	}

	for _, form := range forms {
		index := indexOutsideQuotedText(line, form.markerStart())
		if index < 0 {
			continue
		}

		return form.directive(line[index:]), true
	}

	return "", false
}

func (form commentForm) markerStart() (start string) {
	return form.prefix + markerPrefix
}

func (form commentForm) directive(line string) (directive string) {
	directive = strings.TrimSpace(line[len(form.prefix):])
	if form.suffix != "" {
		directive = strings.TrimSpace(strings.TrimSuffix(directive, form.suffix))
	}

	return directive
}
