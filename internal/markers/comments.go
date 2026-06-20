package markers

import "strings"

// comments constants.
const (
	hashCommentPrefix       = "# "
	slashCommentPrefix      = "// "
	blockCommentPrefix      = "/* "
	blockContinuationPrefix = "* "
	blockCommentSuffix      = "*/"
)

type commentForm struct {
	prefix         string
	optionalSuffix string
}

func extractDirective(line string) (directive string, found bool) {
	forms := [...]commentForm{
		{prefix: hashCommentPrefix},
		{prefix: slashCommentPrefix},
		{prefix: blockCommentPrefix, optionalSuffix: blockCommentSuffix},
		{prefix: blockContinuationPrefix, optionalSuffix: blockCommentSuffix},
	}

	for _, form := range forms {
		index := indexOutsideQuotes(line, form.prefix+markerPrefix)
		if index < 0 {
			continue
		}

		return form.directive(line[index:]), true
	}

	return "", false
}

func (form commentForm) directive(line string) (directive string) {
	directive = strings.TrimSpace(line[len(form.prefix):])
	if form.optionalSuffix != "" {
		directive = strings.TrimSpace(strings.TrimSuffix(directive, form.optionalSuffix))
	}

	return directive
}
