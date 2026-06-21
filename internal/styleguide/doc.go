// Package styleguide parses STYLE.md headings, metadata, and requirement list items.
// It turns Markdown source into a Goldmark AST and returns a document model for coverage checks. It
// does not evaluate rules or decide requirement coverage. Requirement text is normalised plain
// text; Markdown formatting is not preserved.
package styleguide
