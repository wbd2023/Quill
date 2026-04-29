// Package requirementid parses and validates STYLE.md requirement identifiers.
package requirementid

import (
	"fmt"
	"strings"
)

// SectionSlug names IDs written as "<section>.<slug>".
// For example: "3.8.constructor-category-order".
const SectionSlug Scheme = "section_slug"

// Scheme selects a requirement ID grammar.
type Scheme string

// ID is a parsed STYLE.md requirement identifier.
type ID struct {
	value   string
	section string
	slug    string
}

/* ------------------------------------------- Parsing ------------------------------------------ */

// Parse parses a STYLE.md requirement ID according to scheme.
func Parse(value string, scheme Scheme) (id ID, err error) {
	if scheme != SectionSlug {
		return ID{}, fmt.Errorf("unsupported requirement id scheme %q", scheme)
	}

	major, rest, found := strings.Cut(value, ".")
	if !found {
		return ID{}, fmt.Errorf("missing requirement id section")
	}

	minor, slug, found := strings.Cut(rest, ".")
	if !found {
		return ID{}, fmt.Errorf("missing requirement id slug")
	}

	section := major + "." + minor
	if !ValidSection(section) {
		return ID{}, fmt.Errorf("invalid requirement id section %q", section)
	}

	if !isSlug(slug) {
		return ID{}, fmt.Errorf("invalid requirement id slug %q", slug)
	}

	id = ID{
		value:   value,
		section: section,
		slug:    slug,
	}
	return id, nil
}

/* ------------------------------------------ Accessors ----------------------------------------- */

// String returns the original requirement ID.
func (id ID) String() (value string) {
	return id.value
}

// Section returns the numeric STYLE.md section.
func (id ID) Section() (section string) {
	return id.section
}

// Slug returns the readable requirement slug.
func (id ID) Slug() (slug string) {
	return id.slug
}

/* ------------------------------------------ Sections ------------------------------------------ */

// ValidSection reports whether value is a numeric STYLE.md section.
func ValidSection(value string) (valid bool) {
	major, minor, found := strings.Cut(value, ".")
	if !found {
		return false
	}

	return isSectionNumber(major) && isSectionNumber(minor)
}

func isSectionNumber(value string) (valid bool) {
	if value == "" {
		return false
	}

	if value == "0" {
		return true
	}

	if value[0] == '0' {
		return false
	}

	for _, character := range value {
		if character < '0' || character > '9' {
			return false
		}
	}

	return true
}

/* -------------------------------------------- Slugs ------------------------------------------- */

func isSlug(value string) (valid bool) {
	for part := range strings.SplitSeq(value, "-") {
		if !isSlugPart(part) {
			return false
		}
	}

	return true
}

func isSlugPart(value string) (valid bool) {
	if value == "" {
		return false
	}

	for _, character := range value {
		switch {
		case character >= 'a' && character <= 'z':
			continue
		case character >= '0' && character <= '9':
			continue
		default:
			return false
		}
	}

	return true
}
