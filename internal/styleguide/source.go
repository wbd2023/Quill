package styleguide

import (
	"fmt"
	"sort"
)

/* -------------------------------------------- Model ------------------------------------------- */

type sourceFile struct {
	filename string
	contents []byte
	lines    lineTable
}

type position struct {
	line   int
	column int
}

type lineTable []int

type parseError struct {
	filename string
	location position
	message  string
}

/* ---------------------------------------- Source Files ---------------------------------------- */

func newSourceFile(filename string, contents []byte) (file sourceFile) {
	return sourceFile{
		filename: filename,
		contents: contents,
		lines:    newLineTable(contents),
	}
}

/* -------------------------------------- Position Mapping -------------------------------------- */

func (file sourceFile) positionAt(offset int) (location position) {
	return file.lines.positionAt(offset)
}

func newLineTable(contents []byte) (table lineTable) {
	table = make(lineTable, 1)
	for offset, character := range contents {
		if character == '\n' {
			table = append(table, offset+1)
		}
	}

	return table
}

func (table lineTable) positionAt(offset int) (location position) {
	if offset < 0 || len(table) == 0 {
		return position{}
	}

	line := max(sort.Search(len(table), func(index int) bool {
		return table[index] > offset
	})-1, 0)

	return position{
		line:   line + 1,
		column: offset - table[line] + 1,
	}
}

/* ------------------------------------------- Errors ------------------------------------------- */

func (file sourceFile) errorf(location position, format string, arguments ...any) (err error) {
	return parseError{
		filename: file.filename,
		location: location,
		message:  fmt.Sprintf(format, arguments...),
	}
}

func (e parseError) Error() (message string) {
	if e.location.line <= 0 || e.location.column <= 0 {
		return fmt.Sprintf("%s: %s", e.filename, e.message)
	}

	return fmt.Sprintf(
		"%s:%d:%d: %s",
		e.filename,
		e.location.line,
		e.location.column,
		e.message,
	)
}
