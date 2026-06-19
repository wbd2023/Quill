package styleguide

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/style"
)

/* --------------------------------------- Compiler State --------------------------------------- */

type documentCompiler struct {
	file             sourceFile
	scheme           style.IDScheme
	document         Document
	activeSection    string
	pendingMetadata  *requirementMetadata
	seenHeadings     map[string]bool
	seenRequirements map[string]bool
}

func newDocumentCompiler(file sourceFile, scheme style.IDScheme) (compiler documentCompiler) {
	return documentCompiler{
		file:   file,
		scheme: scheme,
		document: Document{
			Headings:     make([]Heading, 0),
			Requirements: make([]Requirement, 0),
		},
		seenHeadings:     make(map[string]bool),
		seenRequirements: make(map[string]bool),
	}
}

/* ----------------------------------------- Compilation ---------------------------------------- */

func (c *documentCompiler) compile(events []documentEvent) (document Document, err error) {
	for _, event := range events {
		if err := c.enterEvent(event); err != nil {
			return Document{}, err
		}
	}

	return c.finish()
}

func (c *documentCompiler) enterEvent(event documentEvent) (err error) {
	switch event.kind {
	case eventHeading:
		return c.enterHeading(event.heading, event.location)

	case eventHTMLBlock:
		return c.enterHTMLBlock(event.text, event.location)

	case eventListItem:
		return c.enterListItem(event.text)

	case eventBoundary:
		return c.enterBoundary()

	default:
		return fmt.Errorf("unknown styleguide document event %d", event.kind)
	}
}

func (c *documentCompiler) finish() (document Document, err error) {
	if err := c.enterBoundary(); err != nil {
		return Document{}, err
	}

	return c.document, nil
}

/* --------------------------------------- Event Handling --------------------------------------- */

func (c *documentCompiler) enterHeading(heading Heading, location position) (err error) {
	if err := c.enterBoundary(); err != nil {
		return err
	}

	return c.recordHeading(heading, location)
}

func (c *documentCompiler) enterHTMLBlock(text string, location position) (err error) {
	fields, found, err := parseMetadataComment(text)
	if err != nil {
		return c.file.errorf(location, "%v", err)
	}

	if !found {
		return c.enterBoundary()
	}

	metadata, err := buildRequirementMetadata(fields, c.scheme)
	if err != nil {
		return c.file.errorf(location, "%v", err)
	}

	metadata.source = location
	return c.recordMetadata(metadata)
}

func (c *documentCompiler) enterListItem(text string) (err error) {
	metadata := c.pendingMetadata
	if metadata == nil {
		return nil
	}

	return c.recordRequirement(*metadata, text)
}

func (c *documentCompiler) enterBoundary() (err error) {
	if c.pendingMetadata != nil {
		return c.unmatchedMetadataError(*c.pendingMetadata)
	}

	return nil
}

/* -------------------------------------- Document Assembly ------------------------------------- */

func (c *documentCompiler) recordHeading(heading Heading, location position) (err error) {
	section := heading.Section
	if c.seenHeadings[section] {
		return c.file.errorf(
			location,
			"duplicate STYLE.md section heading %q",
			section,
		)
	}

	c.activeSection = section
	c.seenHeadings[section] = true

	c.document.Headings = append(c.document.Headings, heading)
	return nil
}

func (c *documentCompiler) recordMetadata(metadata requirementMetadata) (err error) {
	if c.pendingMetadata != nil {
		return c.file.errorf(
			metadata.source,
			"STYLE.md metadata for %q appears before metadata for %q has a requirement list item",
			metadata.id.String(),
			c.pendingMetadata.id.String(),
		)
	}

	c.pendingMetadata = &metadata
	return nil
}

func (c *documentCompiler) recordRequirement(
	metadata requirementMetadata,
	text string,
) (err error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return c.unmatchedMetadataError(metadata)
	}

	id := metadata.id.String()

	activeSection := c.activeSection
	if activeSection == "" {
		return c.file.errorf(
			metadata.source,
			"requirement %q appears before any STYLE.md section heading",
			id,
		)
	}

	expectedSection := metadata.id.Section()
	if expectedSection != activeSection {
		return c.file.errorf(
			metadata.source,
			"requirement %q belongs to section %q but appears under section %q",
			id,
			expectedSection,
			activeSection,
		)
	}

	if c.seenRequirements[id] {
		return c.file.errorf(
			metadata.source,
			"duplicate STYLE.md requirement %q",
			id,
		)
	}

	requirement := Requirement{
		ID:      id,
		Section: activeSection,
		Text:    text,
		Review:  metadata.review,
	}

	c.seenRequirements[id] = true
	c.pendingMetadata = nil

	c.document.Requirements = append(c.document.Requirements, requirement)
	return nil
}

/* ------------------------------------------- Errors ------------------------------------------- */

func (c *documentCompiler) unmatchedMetadataError(metadata requirementMetadata) (err error) {
	return c.file.errorf(
		metadata.source,
		"STYLE.md metadata for %q must be followed by a requirement list item",
		metadata.id.String(),
	)
}
