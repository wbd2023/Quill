package styleguide

import (
	gast "github.com/yuin/goldmark/ast"
)

/* ------------------------------------------- Scanner ------------------------------------------ */

type markdownScanner struct {
	file   sourceFile
	events []documentEvent
}

func newMarkdownScanner(file sourceFile) (scanner markdownScanner) {
	return markdownScanner{
		file:   file,
		events: make([]documentEvent, 0),
	}
}

/* ------------------------------------------ Scanning ------------------------------------------ */

func scanMarkdown(tree gast.Node, file sourceFile) (events []documentEvent) {
	scanner := newMarkdownScanner(file)
	scanner.scan(tree)
	return scanner.events
}

func (s *markdownScanner) scan(tree gast.Node) {
	_ = gast.Walk(tree, func(node gast.Node, entering bool) (status gast.WalkStatus, err error) {
		if !entering {
			return gast.WalkContinue, nil
		}

		s.scanNode(node)
		return gast.WalkContinue, nil
	})
}

func (s *markdownScanner) scanNode(node gast.Node) {
	switch node := node.(type) {
	case *gast.Heading:
		s.emitHeading(node)

	case *gast.HTMLBlock:
		s.emitHTMLBlock(node)

	case *gast.ListItem:
		s.emitListItem(node)

	default:
		s.emitBoundary(node)
	}
}

/* --------------------------------------- Event Emission --------------------------------------- */

func (s *markdownScanner) emitHeading(node *gast.Heading) {
	location := s.locationOf(node)
	heading, found := parseHeading(extractInlineText(s.file.contents, node))
	if !found {
		s.events = append(s.events, newBoundaryEvent(location))
		return
	}

	s.events = append(s.events, newHeadingEvent(location, heading))
}

func (s *markdownScanner) emitHTMLBlock(node *gast.HTMLBlock) {
	s.events = append(s.events, newHTMLBlockEvent(
		s.locationOf(node),
		extractHTMLBlockText(s.file.contents, node),
	))
}

func (s *markdownScanner) emitListItem(node *gast.ListItem) {
	s.events = append(s.events, newListItemEvent(
		s.locationOf(node),
		extractInlineText(s.file.contents, node),
	))
}

func (s *markdownScanner) emitBoundary(node gast.Node) {
	if !isBoundaryNode(node) {
		return
	}

	s.events = append(s.events, newBoundaryEvent(s.locationOf(node)))
}

/* --------------------------------------- Boundary Nodes --------------------------------------- */

// Boundary nodes prevent pending metadata from attaching to later list items.
func isBoundaryNode(node gast.Node) (boundary bool) {
	if node.Type() != gast.TypeBlock {
		return false
	}

	if hasListItemAncestor(node) {
		return false
	}

	switch node.(type) {
	case *gast.Document, *gast.List:
		return false
	case *gast.Heading, *gast.HTMLBlock, *gast.ListItem:
		return false
	default:
		return true
	}
}

func hasListItemAncestor(node gast.Node) (found bool) {
	for parent := node.Parent(); parent != nil; parent = parent.Parent() {
		if _, ok := parent.(*gast.ListItem); ok {
			return true
		}
	}

	return false
}

/* -------------------------------------- Source Locations -------------------------------------- */

func (s *markdownScanner) locationOf(node gast.Node) (location position) {
	offset := node.Pos()
	if offset < 0 && node.Type() == gast.TypeBlock {
		// Goldmark reports Pos as -1 for some block nodes; their source span lives in Lines.
		lines := node.Lines()
		if lines.Len() > 0 {
			offset = lines.At(0).Start
		}
	}

	return s.file.positionAt(offset)
}
