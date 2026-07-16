package style

import (
	"fmt"
	"strings"
)

// RequirementID is a parsed STYLE.md requirement identifier.
type RequirementID struct {
	value   string
	section string
	slug    string
}

/* ------------------------------------------- Parsing ------------------------------------------ */

// ParseRequirementID parses a STYLE.md requirement identifier in "<section>.<slug>" form.
func ParseRequirementID(value string) (id RequirementID, err error) {
	major, rest, found := strings.Cut(value, ".")
	if !found {
		return RequirementID{}, fmt.Errorf("missing requirement id section")
	}

	minor, slug, found := strings.Cut(rest, ".")
	if !found {
		return RequirementID{}, fmt.Errorf("missing requirement id slug")
	}

	section := major + "." + minor
	if !IsValidSection(section) {
		return RequirementID{}, fmt.Errorf("invalid requirement id section %q", section)
	}

	if !isSlug(slug) {
		return RequirementID{}, fmt.Errorf("invalid requirement id slug %q", slug)
	}

	id = RequirementID{
		value:   value,
		section: section,
		slug:    slug,
	}
	return id, nil
}

/* ------------------------------------------ Accessors ----------------------------------------- */

// String returns the original requirement ID.
func (id RequirementID) String() (value string) {
	return id.value
}

// Section returns the numeric STYLE.md section.
func (id RequirementID) Section() (section string) {
	return id.section
}

// Slug returns the readable requirement slug.
func (id RequirementID) Slug() (slug string) {
	return id.slug
}

/* ------------------------------------------ Sections ------------------------------------------ */

// IsValidSection reports whether value is a numeric STYLE.md section.
func IsValidSection(value string) (valid bool) {
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
