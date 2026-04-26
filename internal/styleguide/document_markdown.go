package styleguide

import (
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

/* --------------------------------------- Heading Parsing -------------------------------------- */

func parseHeading(line string) (section string, title string, found bool) {
	return parseHeadingText(line)
}

/* ------------------------------------- Requirement Parsing ------------------------------------ */

func parseRequirement(line string) (requirementID string, text string, found bool) {
	return parseRequirementText(line, "", RequirementIDFormatSectionSlug)
}

func parseRequirementID(line string) (requirementID string, found bool) {
	requirementID, _, found = parseRequirement(line)
	return requirementID, found
}

/* ------------------------------------ Node Text Extraction ------------------------------------ */

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

/* ---------------------------------------- List Parsing ---------------------------------------- */

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
