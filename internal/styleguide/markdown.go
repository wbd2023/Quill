package styleguide

import (
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

func nodePlainText(source []byte, node gast.Node) (plain string) {
	var builder strings.Builder
	_ = gast.Walk(node, func(child gast.Node, entering bool) (status gast.WalkStatus, err error) {
		if !entering {
			return gast.WalkContinue, nil
		}

		switch typedNode := child.(type) {
		case *gast.Text:
			builder.Write(typedNode.Segment.Value(source))
			if typedNode.HardLineBreak() || typedNode.SoftLineBreak() {
				builder.WriteByte(' ')
			}

		case *gast.String:
			builder.Write(typedNode.Value)
		}

		return gast.WalkContinue, nil
	})

	return strings.Join(strings.Fields(builder.String()), " ")
}

func htmlBlockText(source []byte, node *gast.HTMLBlock) (text string) {
	var builder strings.Builder
	for lineIndex := 0; lineIndex < node.Lines().Len(); lineIndex++ {
		segment := node.Lines().At(lineIndex)
		builder.Write(segment.Value(source))
	}
	if node.HasClosure() {
		builder.Write(node.ClosureLine.Value(source))
	}

	return builder.String()
}

func parseListItemBody(line string) (body string, found bool) {
	trimmed := strings.TrimSpace(line)
	switch {
	case strings.HasPrefix(trimmed, "* "):
		return strings.TrimSpace(trimmed[2:]), true

	case startsWithOrderedListMarker(trimmed):
		markerEnd := strings.Index(trimmed, ". ")
		return strings.TrimSpace(trimmed[markerEnd+2:]), true

	default:
		return "", false
	}
}

func startsWithOrderedListMarker(value string) (found bool) {
	dotIndex := strings.Index(value, ". ")
	if dotIndex <= 0 {
		return false
	}

	for _, character := range value[:dotIndex] {
		if character < '0' || character > '9' {
			return false
		}
	}

	return true
}
