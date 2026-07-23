package text

import (
	"bytes"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	gmtext "github.com/yuin/goldmark/text"
)

/* ------------------------------------- Markdown allowance ------------------------------------- */

// markdownAllowances maps a 1-based source line number to the verbatim destination
// bytes of the single qualifying link reference definition whose destination lives
// on that line. A nil map means no line receives the overlong allowance.
type markdownAllowances map[int][]byte

// loadMarkdownAllowances parses path once when it is Markdown and returns the set
// of lines whose overlong status the narrow exception may excuse. Any read or
// parse problem resolves to a nil map so the caller fails closed (flags the line).
func loadMarkdownAllowances(path string) (allowances markdownAllowances) {
	if !isMarkdownFile(path) {
		return nil
	}

	source, err := os.ReadFile(path)
	if err != nil {
		// ScanLines reports the genuine read error; here we fail closed for allowance.
		return nil
	}

	return buildMarkdownAllowances(source)
}

// allows reports whether an overlong line is excused by a qualifying destination.
// It verifies the destination appears exactly once in the line text, the
// destination's own tab-expanded width strictly exceeds the limit, and the
// residual after removing only the destination bytes stays within the limit.
func (a markdownAllowances) allows(lineNumber int, line string) (allowed bool) {
	destination, ok := a[lineNumber]
	if !ok {
		return false
	}

	needle := string(destination)
	index := strings.Index(line, needle)
	if index < 0 {
		return false
	}

	if strings.Index(line[index+len(needle):], needle) >= 0 {
		// Ambiguous layout: cannot pin the destination bytes. Fail closed.
		return false
	}

	if len(expandTabs(needle)) <= lineLengthLimit {
		return false
	}

	residual := line[:index] + line[index+len(needle):]
	if len(expandTabs(residual)) > lineLengthLimit {
		return false
	}

	return true
}

/* --------------------------------------- Allowance build -------------------------------------- */

// buildMarkdownAllowances parses source with the default Goldmark parser and
// records lowercase http(s) link reference definition destinations keyed by line.
func buildMarkdownAllowances(source []byte) (allowances markdownAllowances) {
	document := goldmark.DefaultParser().Parse(gmtext.NewReader(source))
	lines := newLineIndex(source)
	allowances = markdownAllowances{}

	_ = gast.Walk(document, func(node gast.Node, entering bool) (gast.WalkStatus, error) {
		if !entering {
			return gast.WalkContinue, nil
		}

		definition, ok := node.(*gast.LinkReferenceDefinition)
		if !ok {
			return gast.WalkContinue, nil
		}

		destination := definition.Destination
		if !isLowercaseHTTPDestination(destination) {
			return gast.WalkContinue, nil
		}

		if definition.Lines().Len() == 0 {
			return gast.WalkContinue, nil
		}

		start, ok := destinationSourceOffset(source, definition.Lines(), destination)
		if !ok {
			return gast.WalkContinue, nil
		}

		lineNumber := lines.lineNumber(start)
		if _, exists := allowances[lineNumber]; exists {
			// Two definitions cannot share a first line; treat as ambiguous and fail closed.
			delete(allowances, lineNumber)
			return gast.WalkContinue, nil
		}

		allowances[lineNumber] = destination
		return gast.WalkContinue, nil
	})

	return allowances
}

// destinationSourceOffset locates the real source byte offset of a link reference
// definition destination via a structural cursor over every source segment the
// parser recorded for the definition. The previous implementation searched only
// the first segment, which both missed a destination living on a continuation
// line (after the permitted post-colon line break) and could match the same
// bytes inside the label.
//
// The cursor reconstructs the parser's logical view of the definition - label,
// colon, the whitespace the spec permits between them and the destination
// (including the line break that lets the destination begin on a later line) -
// from the concatenated segment values, and accepts the destination only when
// its exact raw bytes occupy that single syntax position. Each segment's Padding
// virtual spaces are stripped leading indentation with no physical source byte,
// so they map to no offset and never receive an allowance. Any absence,
// ambiguity, or offset mismatch fails closed.
func destinationSourceOffset(
	source []byte,
	segments *gmtext.Segments,
	destination []byte,
) (offset int, ok bool) {
	if segments == nil || segments.Len() == 0 || len(destination) == 0 {
		return 0, false
	}

	// logical holds the parser's view of the definition text; srcOf maps each
	// logical byte back to its physical source offset, or -1 for Padding virtual
	// spaces that have no source counterpart. Stripped inter-segment indentation
	// (e.g. list-item content offsets) is absent from the segments and therefore
	// absent here too.
	logical := make([]byte, 0, segments.Len()*32)
	srcOf := make([]int, 0, cap(logical))
	for index := range segments.Len() {
		segment := segments.At(index)
		for range segment.Padding {
			logical = append(logical, ' ')
			srcOf = append(srcOf, -1)
		}
		for sourceIndex := segment.Start; sourceIndex < segment.Stop; sourceIndex++ {
			logical = append(logical, source[sourceIndex])
			srcOf = append(srcOf, sourceIndex)
		}
	}

	cursor := skipDefinitionSpaces(logical, 0)
	if cursor >= len(logical) || logical[cursor] != '[' {
		return 0, false
	}

	// Advance past the label to its closing ']'. The closer is the first ']'
	// that is not backslash-escaped.
	cursor++
	closeAt := -1
	for cursor < len(logical) {
		switch {
		case logical[cursor] == '\\' && cursor+1 < len(logical):
			cursor += 2
		case logical[cursor] == ']':
			closeAt = cursor
		default:
			cursor++
		}
		if closeAt >= 0 {
			break
		}
	}
	if closeAt < 0 {
		return 0, false
	}

	// Require the colon immediately after the label closer.
	if closeAt+1 >= len(logical) || logical[closeAt+1] != ':' {
		return 0, false
	}

	// Consume the whitespace the spec permits between the colon and the
	// destination: spaces, tabs, and the line break that lets the destination
	// begin on a continuation line.
	cursor = skipDefinitionSpaces(logical, closeAt+2)
	if cursor >= len(logical) {
		return 0, false
	}

	// A destination may be wrapped in optional angle brackets. Goldmark
	// strips the delimiters from the parsed value, so recognize a leading
	// '<' here and require the matching '>' immediately after the
	// destination bytes. Only the destination's own physical bytes are
	// mapped to an offset - the delimiters are no part of the allowance.
	destStart := cursor
	needCloser := false
	if logical[destStart] == '<' {
		destStart++
		needCloser = true
	}

	if destStart+len(destination) > len(logical) {
		return 0, false
	}

	if !bytes.Equal(logical[destStart:destStart+len(destination)], destination) {
		return 0, false
	}
	if needCloser {
		closer := destStart + len(destination)
		if closer >= len(logical) || logical[closer] != '>' {
			return 0, false
		}
	}

	physical := srcOf[destStart]
	if physical < 0 || physical+len(destination) > len(source) {
		return 0, false
	}

	if !bytes.Equal(source[physical:physical+len(destination)], destination) {
		return 0, false
	}

	return physical, true
}

// skipDefinitionSpaces advances cursor past spaces, tabs, and line breaks in the
// reconstructed definition text, stopping at the first non-whitespace byte.
func skipDefinitionSpaces(logical []byte, cursor int) (position int) {
	for cursor < len(logical) {
		switch logical[cursor] {
		case ' ', '\t', '\n', '\r':
			cursor++
		default:
			return cursor
		}
	}
	return cursor
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// isMarkdownFile reports whether path is a Markdown source file the exception covers.
func isMarkdownFile(path string) (markdown bool) {
	return strings.EqualFold(filepath.Ext(path), ".md")
}

// isLowercaseHTTPDestination reports whether destination begins with an exact
// lowercase http:// or https:// scheme.
func isLowercaseHTTPDestination(destination []byte) (lowercase bool) {
	return bytes.HasPrefix(destination, []byte("http://")) ||
		bytes.HasPrefix(destination, []byte("https://"))
}

// expandTabs replaces tabs with tab-width spaces, matching the line-length count.
func expandTabs(value string) (expanded string) {
	return strings.ReplaceAll(value, "\t", strings.Repeat(" ", lineLengthTabWidth))
}

// lineIndex holds the source byte offset at which each 1-based line begins.
type lineIndex []int

func newLineIndex(source []byte) (index lineIndex) {
	index = append(index, 0)
	for offset, character := range source {
		if character == '\n' {
			index = append(index, offset+1)
		}
	}
	return index
}

func (index lineIndex) lineNumber(offset int) (number int) {
	if offset < 0 {
		return 1
	}
	firstGreater := sort.Search(len(index), func(line int) bool {
		return index[line] > offset
	})
	if firstGreater == 0 {
		return 1
	}
	return firstGreater
}
