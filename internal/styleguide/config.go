package styleguide

const defaultFilename = "STYLE.md"

// Config controls STYLE.md parsing.
type Config struct {
	// Filename names the STYLE.md file relative to the repository root. Load requires it; Parse
	// uses it only in diagnostics and defaults to STYLE.md.
	Filename string
}
