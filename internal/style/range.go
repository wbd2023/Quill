package style

import (
	"errors"
	"fmt"
	"path"
	"strings"
)

/* ----------------------------------------- Range Model ---------------------------------------- */

// Position is a one-based UTF-8 byte location within a file. Line is one-based; Column is the
// one-based byte offset within Line. A zero Line means the position is entirely unknown; a zero
// Column with a non-zero Line means the line is known but the column is not. Columns count UTF-8
// bytes, not runes, so a multi-byte character occupies as many columns as it has bytes.
type Position struct {
	Line   int
	Column int
}

// IsKnown reports whether the position carries a usable line. A position is known once its line is
// set; the column may still be unknown (zero).
func (position Position) IsKnown() (known bool) {
	return position.Line > 0
}

// Range is a half-open span within a file: Start is inclusive and End is exclusive. The zero Range
// represents an unknown location. Either end may be unknown independently: a known Start with an
// unknown End means the location begins at Start and its extent is unspecified.
type Range struct {
	Start Position
	End   Position
}

// IsUnknown reports whether the range carries no usable location at all. A range is unknown only
// when neither end has a known line.
func (r Range) IsUnknown() (unknown bool) {
	return !r.Start.IsKnown() && !r.End.IsKnown()
}

/* -------------------------------------- Boundary Verifier ------------------------------------- */

// VerifyRange is the diagnostic protocol-boundary check. It rejects, rather than repairs, file
// paths and ranges that are not safe to admit into Quill's trusted diagnostic model: file must be
// a clean repository-relative slash path with no escape, positions must be valid one-based UTF-8
// byte locations, and a known End must not precede its Start. External Pack diagnostics must pass
// this check at the point they cross into Quill; built-in producers construct valid ranges
// directly.
func VerifyRange(file string, location Range) (err error) {
	if err = verifyFile(file); err != nil {
		return err
	}

	// A diagnostic without a file is repository-level and must carry no location: a position
	// without a file is meaningless and would render as a dangling ":line".
	if file == "" && (location.Start.IsKnown() || location.End.IsKnown()) {
		return errRangeWithoutFile
	}

	if err = verifyPosition(location.Start, "start"); err != nil {
		return err
	}
	if err = verifyPosition(location.End, "end"); err != nil {
		return err
	}

	// An End without a Start is meaningless and therefore rejected.
	if location.End.IsKnown() && !location.Start.IsKnown() {
		return errRangeEndWithoutStart
	}

	// When both ends are known, End must not precede Start. Columns are compared only when both
	// ends carry a column, since a zero column means "column unknown".
	if location.Start.IsKnown() && location.End.IsKnown() {
		if location.End.Line < location.Start.Line {
			return errRangeEndBeforeStart
		}
		if location.End.Line == location.Start.Line &&
			location.Start.Column > 0 && location.End.Column > 0 &&
			location.End.Column < location.Start.Column {
			return errRangeEndBeforeStart
		}
	}

	return nil
}

// verifyFile rejects paths that are not clean, repository-relative, forward-slash paths. It
// accepts an empty file (a repository- or Pack-level finding with no single file).
func verifyFile(file string) (err error) {
	if file == "" {
		return nil
	}

	if strings.ContainsAny(file, "\\\x00") {
		return errPathUnclean
	}

	if path.IsAbs(file) {
		return errPathAbsolute
	}

	// A leading Windows drive prefix (for example C:/x or C:\x) is an absolute path in disguise.
	if len(file) >= 2 && isASCIILetter(file[0]) && file[1] == ':' {
		return errPathAbsolute
	}

	for _, segment := range strings.Split(file, "/") {
		if segment == "" || segment == "." || segment == ".." {
			return errPathEscape
		}
	}

	return nil
}

// verifyPosition rejects negative coordinates and a column set without a line.
func verifyPosition(position Position, label string) (err error) {
	if position.Line < 0 || position.Column < 0 {
		return fmt.Errorf("%w: %s", errPositionNegative, label)
	}
	if position.Line == 0 && position.Column > 0 {
		return fmt.Errorf("%w: %s", errColumnWithoutLine, label)
	}

	return nil
}

func isASCIILetter(byteValue byte) (letter bool) {
	return (byteValue >= 'a' && byteValue <= 'z') || (byteValue >= 'A' && byteValue <= 'Z')
}

var (
	errPathAbsolute         = errors.New("diagnostic file path must be repository-relative")
	errPathEscape           = errors.New("diagnostic file path must not escape the repository")
	errPathUnclean          = errors.New("diagnostic file path must be a clean slash path")
	errColumnWithoutLine    = errors.New("diagnostic column set without a line")
	errPositionNegative     = errors.New("diagnostic position must be non-negative")
	errRangeEndBeforeStart  = errors.New("diagnostic range end must not precede start")
	errRangeEndWithoutStart = errors.New("diagnostic range end set without a start")
	errRangeWithoutFile     = errors.New("diagnostic position set without a file")
)
