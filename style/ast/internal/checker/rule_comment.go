package checker

import (
	"go/ast"
	"go/token"
	"strings"

	"stylecheck/internal/checker/support"
)

const inlineCommentDirectiveCodeGenerated = "code generated"
const inlineCommentDirectiveFixme = "fixme:"
const inlineCommentDirectiveGo = "go:"
const inlineCommentDirectiveNolint = "nolint"
const inlineCommentDirectiveTodo = "todo:"

// checkInlineCommentStyle validates trailing inline comment case and punctuation (2.3).
func checkInlineCommentStyle(
	fileSet *token.FileSet,
	file *ast.File,
	path string,
) (violations []violation) {
	if !support.IsAppScopePath(path) {
		return nil
	}

	commentMap := ast.NewCommentMap(fileSet, file, file.Comments)
	seen := make(map[token.Pos]bool)

	for node, commentGroups := range commentMap {
		nodeEndLine := fileSet.Position(node.End()).Line

		for _, commentGroup := range commentGroups {
			for _, comment := range commentGroup.List {
				if !strings.HasPrefix(comment.Text, "//") {
					continue
				}

				if seen[comment.Pos()] {
					continue
				}

				commentPosition := fileSet.Position(comment.Pos())
				if commentPosition.Line != nodeEndLine {
					continue
				}

				payload := strings.TrimSpace(strings.TrimPrefix(comment.Text, "//"))
				if payload == "" || isInlineCommentDirective(payload) {
					continue
				}

				seen[comment.Pos()] = true

				if support.StartsWithUppercaseLetter(payload) {
					violations = append(violations, violation{
						position: fileSet.Position(comment.Pos()),
						rule:     "2.3",
						message:  "inline trailing comment should start lower-case",
					})
				}

				if support.EndsWithSentencePunctuation(payload) {
					violations = append(violations, violation{
						position: fileSet.Position(comment.Pos()),
						rule:     "2.3",
						message:  "inline trailing comment should not end with punctuation",
					})
				}
			}
		}
	}

	return violations
}

func isInlineCommentDirective(comment string) (found bool) {
	normalisedComment := strings.ToLower(strings.TrimSpace(comment))

	return strings.HasPrefix(normalisedComment, inlineCommentDirectiveNolint) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveTodo) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveFixme) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveGo) ||
		strings.HasPrefix(normalisedComment, inlineCommentDirectiveCodeGenerated)
}
