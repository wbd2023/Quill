package styleguide

import (
	"fmt"
	"strings"

	gast "github.com/yuin/goldmark/ast"
)

/* -------------------------------------------- Types ------------------------------------------- */

type documentWalkState struct {
	document            documentModel
	currentSection      string
	pendingMetadata     *RequirementMetadata
	requirementIDFormat string
	seenHeadings        map[string]bool
	seenRequirements    map[string]bool
}

/* ----------------------------------- Walk State Construction ---------------------------------- */

func newDocumentWalkState(requirementIDFormat string) (state documentWalkState) {
	return documentWalkState{
		document: documentModel{
			Headings:     make([]documentHeading, 0),
			Requirements: make([]documentRequirement, 0),
		},
		requirementIDFormat: requirementIDFormat,
		seenHeadings:        make(map[string]bool),
		seenRequirements:    make(map[string]bool),
	}
}

/* -------------------------------------- Document Walking -------------------------------------- */

func (state *documentWalkState) walk(
	source []byte,
) (walker func(gast.Node, bool) (status gast.WalkStatus, err error)) {
	return func(node gast.Node, entering bool) (status gast.WalkStatus, err error) {
		if !entering {
			return gast.WalkContinue, nil
		}

		switch typedNode := node.(type) {
		case *gast.Heading:
			return state.enterHeading(source, typedNode)

		case *gast.HTMLBlock:
			return state.enterHTMLBlock(source, typedNode)

		case *gast.ListItem:
			return state.enterListItem(source, typedNode)

		default:
			return gast.WalkContinue, nil
		}
	}
}

/* -------------------------------------- Heading Handling -------------------------------------- */

func (state *documentWalkState) enterHeading(
	source []byte,
	node *gast.Heading,
) (status gast.WalkStatus, err error) {
	if state.pendingMetadata != nil {
		return gast.WalkStop, fmt.Errorf(
			"style.md metadata for %q must be followed by a requirement bullet",
			state.pendingMetadata.ID,
		)
	}

	section, title, found := parseHeadingText(nodePlainText(source, node))
	if !found {
		return gast.WalkContinue, nil
	}

	if state.seenHeadings[section] {
		return gast.WalkStop, fmt.Errorf("duplicate style.md section heading %q", section)
	}

	state.currentSection = section
	state.seenHeadings[section] = true
	state.document.Headings = append(state.document.Headings, documentHeading{
		Section: section,
		Title:   title,
	})
	return gast.WalkContinue, nil
}

/* -------------------------------------- Metadata Handling ------------------------------------- */

func (state *documentWalkState) enterHTMLBlock(
	source []byte,
	node *gast.HTMLBlock,
) (status gast.WalkStatus, err error) {
	metadata, found, err := parseMetadataComment(
		htmlBlockText(source, node),
		state.requirementIDFormat,
	)
	if err != nil {
		return gast.WalkStop, err
	}

	if !found {
		if state.pendingMetadata != nil {
			return gast.WalkStop, fmt.Errorf(
				"style.md metadata for %q must be followed by a requirement bullet",
				state.pendingMetadata.ID,
			)
		}

		return gast.WalkContinue, nil
	}

	if state.pendingMetadata != nil {
		return gast.WalkStop, fmt.Errorf(
			"style.md metadata for %q must be followed by a requirement before another "+
				"metadata comment",
			state.pendingMetadata.ID,
		)
	}

	state.pendingMetadata = &metadata
	return gast.WalkContinue, nil
}

func htmlBlockText(source []byte, node *gast.HTMLBlock) (text string) {
	var builder strings.Builder
	for index := 0; index < node.Lines().Len(); index++ {
		segment := node.Lines().At(index)
		builder.Write(segment.Value(source))
	}
	if node.HasClosure() {
		builder.Write(node.ClosureLine.Value(source))
	}

	return builder.String()
}

/* ------------------------------------ Requirement Handling ------------------------------------ */

func (state *documentWalkState) enterListItem(
	source []byte,
	node *gast.ListItem,
) (status gast.WalkStatus, err error) {
	pendingRequirementID := ""
	if state.pendingMetadata != nil {
		pendingRequirementID = state.pendingMetadata.ID
	}

	requirementID, requirementText, found := parseRequirementText(
		nodePlainText(source, node),
		pendingRequirementID,
		state.requirementIDFormat,
	)
	if !found {
		if state.pendingMetadata != nil {
			return gast.WalkStop, fmt.Errorf(
				"style.md metadata for %q must be followed by a requirement bullet",
				state.pendingMetadata.ID,
			)
		}

		return gast.WalkContinue, nil
	}

	if state.currentSection == "" {
		return gast.WalkStop, fmt.Errorf(
			"requirement %q appears before any style.md section heading",
			requirementID,
		)
	}

	if requirementSection(requirementID, state.requirementIDFormat) != state.currentSection {
		return gast.WalkStop, fmt.Errorf(
			"requirement %q appears under section %q",
			requirementID,
			state.currentSection,
		)
	}

	if state.seenRequirements[requirementID] {
		return gast.WalkStop, fmt.Errorf("duplicate style.md requirement %q", requirementID)
	}

	state.seenRequirements[requirementID] = true
	requirement := documentRequirement{
		ID:      requirementID,
		Section: state.currentSection,
		Text:    requirementText,
	}
	if state.pendingMetadata != nil {
		requirement.Mode = state.pendingMetadata.Mode
		requirement.Reason = state.pendingMetadata.Reason
		state.pendingMetadata = nil
	}

	state.document.Requirements = append(state.document.Requirements, requirement)
	return gast.WalkContinue, nil
}

/* -------------------------------------- Walk Finalisation ------------------------------------- */

func (state *documentWalkState) finish() (document documentModel, err error) {
	if state.pendingMetadata != nil {
		return documentModel{}, fmt.Errorf(
			"style.md metadata for %q must be followed by a requirement bullet",
			state.pendingMetadata.ID,
		)
	}

	return state.document, nil
}
