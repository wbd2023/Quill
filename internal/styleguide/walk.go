package styleguide

import (
	"fmt"

	gast "github.com/yuin/goldmark/ast"
)

/* ----------------------------------------- Walk State ----------------------------------------- */

type documentWalkState struct {
	document            Document
	activeSection       string
	pendingMetadata     *RequirementMetadata
	requirementIDFormat string
	seenHeadings        map[string]bool
	seenRequirements    map[string]bool
}

func newDocumentWalkState(requirementIDFormat string) (state documentWalkState) {
	return documentWalkState{
		document: Document{
			Headings:     make([]Heading, 0),
			Requirements: make([]Requirement, 0),
		},
		requirementIDFormat: requirementIDFormat,
		seenHeadings:        make(map[string]bool),
		seenRequirements:    make(map[string]bool),
	}
}

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

func (state *documentWalkState) finish() (document Document, err error) {
	if state.pendingMetadata != nil {
		return Document{}, state.pendingMetadataError("requirement bullet")
	}

	return state.document, nil
}

/* -------------------------------------- Heading Handling -------------------------------------- */

func (state *documentWalkState) enterHeading(
	source []byte,
	node *gast.Heading,
) (status gast.WalkStatus, err error) {
	if state.pendingMetadata != nil {
		return gast.WalkStop, state.pendingMetadataError("requirement bullet")
	}

	section, title, found := parseHeadingText(nodePlainText(source, node))
	if !found {
		return gast.WalkContinue, nil
	}

	if state.seenHeadings[section] {
		return gast.WalkStop, fmt.Errorf("duplicate style.md section heading %q", section)
	}

	state.activeSection = section
	state.seenHeadings[section] = true
	state.document.Headings = append(state.document.Headings, Heading{
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
			return gast.WalkStop, state.pendingMetadataError("requirement bullet")
		}

		return gast.WalkContinue, nil
	}

	if state.pendingMetadata != nil {
		return gast.WalkStop, state.pendingMetadataError(
			"requirement before another metadata comment",
		)
	}

	state.pendingMetadata = &metadata
	return gast.WalkContinue, nil
}

func (state *documentWalkState) pendingMetadataError(expectation string) (err error) {
	return fmt.Errorf(
		"style.md metadata for %q must be followed by a %s",
		state.pendingMetadata.ID,
		expectation,
	)
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
			return gast.WalkStop, state.pendingMetadataError("requirement bullet")
		}

		return gast.WalkContinue, nil
	}

	if err := state.validateRequirementLocation(requirementID); err != nil {
		return gast.WalkStop, err
	}

	state.seenRequirements[requirementID] = true
	requirement := Requirement{
		ID:      requirementID,
		Section: state.activeSection,
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

func (state *documentWalkState) validateRequirementLocation(requirementID string) (err error) {
	if state.activeSection == "" {
		return fmt.Errorf(
			"requirement %q appears before any style.md section heading",
			requirementID,
		)
	}

	if requirementSection(requirementID, state.requirementIDFormat) != state.activeSection {
		return fmt.Errorf(
			"requirement %q appears under section %q",
			requirementID,
			state.activeSection,
		)
	}

	if state.seenRequirements[requirementID] {
		return fmt.Errorf("duplicate style.md requirement %q", requirementID)
	}

	return nil
}
