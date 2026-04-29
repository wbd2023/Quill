package styleguide

import (
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

func extractInlineText(contents []byte, node gast.Node) (text string) {
	var builder strings.Builder
	_ = gast.Walk(node, func(child gast.Node, entering bool) (status gast.WalkStatus, err error) {
		if !entering {
			return gast.WalkContinue, nil
		}

		// Block nodes mark paragraph/list boundaries without contributing text themselves.
		if child != node && child.Type() == gast.TypeBlock && builder.Len() > 0 {
			builder.WriteByte(' ')
		}

		switch child := child.(type) {
		case *gast.Text:
			builder.Write(child.Segment.Value(contents))
			if child.HardLineBreak() || child.SoftLineBreak() {
				builder.WriteByte(' ')
			}

		case *gast.String:
			builder.Write(child.Value)
		}

		return gast.WalkContinue, nil
	})

	return strings.Join(strings.Fields(builder.String()), " ")
}

func extractHTMLBlockText(contents []byte, node *gast.HTMLBlock) (text string) {
	value := node.Lines().Value(contents)
	if node.HasClosure() {
		value = append(value, node.ClosureLine.Value(contents)...)
	}

	return string(value)
}
