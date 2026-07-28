package external

import (
	"bufio"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	lineScannerStartBuffer = 64 * 1024
	lineScannerMaxBuffer   = 1024 * 1024
)

// lineScanner reads JSON Lines, decoding each non-blank line into an envelope. Blank lines are
// skipped so a Pack may separate records with whitespace. A line that is not valid JSON becomes a
// scan error surfaced through Err, terminating iteration.
type lineScanner struct {
	lines  *bufio.Scanner
	record envelope
	err    error
}

func newLineScanner(stdout string) (scanner *lineScanner) {
	lines := bufio.NewScanner(strings.NewReader(stdout))
	lines.Buffer(make([]byte, 0, lineScannerStartBuffer), lineScannerMaxBuffer)
	return &lineScanner{lines: lines}
}

// Scan advances to the next non-blank line, decoding it. It returns false when no more records
// remain or when a malformed line is encountered; in the latter case Err reports the cause.
func (scanner *lineScanner) Scan() (hasNext bool) {
	for scanner.lines.Scan() {
		raw := scanner.lines.Text()
		if strings.TrimSpace(raw) == "" {
			continue
		}

		scanner.record = envelope{}
		if err := json.Unmarshal([]byte(raw), &scanner.record); err != nil {
			scanner.err = fmt.Errorf("external pack emitted invalid JSON: %w", err)
			return false
		}
		return true
	}

	return false
}

// Text returns the envelope decoded for the current record.
func (scanner *lineScanner) Text() (record envelope) {
	return scanner.record
}

// Err returns the error that stopped iteration, or nil when the stream ended cleanly.
func (scanner *lineScanner) Err() (err error) {
	if scanner.err != nil {
		return scanner.err
	}
	return scanner.lines.Err()
}
