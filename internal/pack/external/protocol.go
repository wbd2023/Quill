package external

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/style"
)

/* ------------------------------------------- Request ------------------------------------------ */

// Request is the single JSON object Quill writes to an external Pack's standard input.
type Request struct {
	Protocol       string         `json:"protocol"`
	Operation      string         `json:"operation"`
	RepositoryRoot string         `json:"repository_root"`
	PackID         string         `json:"pack_id"`
	RuleID         string         `json:"rule_id"`
	CheckID        string         `json:"check_id"`
	Scope          string         `json:"scope"`
	Files          []string       `json:"files"`
	Configuration  map[string]any `json:"configuration"`
}

// EncodeRequest marshals request as the standard-input payload.
func EncodeRequest(request Request) (payload []byte, err error) {
	encoded, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("encode external pack request: %w", err)
	}
	return encoded, nil
}

/* ------------------------------------------ Response ------------------------------------------ */

// Outcome is the parsed result of one external Pack invocation: the diagnostics emitted before the
// terminal completion, and the completion status.
type Outcome struct {
	Diagnostics []style.Diagnostic
	Success     bool
	// Error carries the Pack-reported failure message when Success is false.
	Error string
}

type envelope struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	File    string `json:"file"`
	Start   *pos   `json:"start"`
	End     *pos   `json:"end"`
	HelpURL string `json:"help_url"`

	Success bool   `json:"success"`
	ErrText string `json:"error"`
}

type pos struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// ParseResponse decodes the JSON Lines an external Pack wrote to standard output. The protocol is
// strict: every non-blank line is one record, records are diagnostics until exactly one terminal
// completion, and every diagnostic path and range is verified at the protocol boundary through
// style.VerifyRange before it is admitted into Quill's trusted diagnostic model.
func ParseResponse(stdout string) (outcome Outcome, err error) {
	scanner := newLineScanner(stdout)

	var line int
	var completed bool
	diagnostics := []style.Diagnostic(nil)

	for scanner.Scan() {
		line++
		record := scanner.Text()

		switch record.Type {
		case "diagnostic":
			if completed {
				return Outcome{}, fmt.Errorf(
					"external pack emitted a diagnostic after completion (line %d)", line)
			}

			diagnostic, diagErr := decodeDiagnostic(record, line)
			if diagErr != nil {
				return Outcome{}, diagErr
			}
			diagnostics = append(diagnostics, diagnostic)

		case "complete":
			if completed {
				return Outcome{}, fmt.Errorf(
					"external pack emitted more than one completion (line %d)", line)
			}
			completed = true
			outcome.Success = record.Success
			outcome.Error = record.ErrText

		default:
			return Outcome{}, fmt.Errorf(
				"external pack emitted record with unknown type %q (line %d)", record.Type, line)
		}
	}

	if scanErr := scanner.Err(); scanErr != nil {
		return Outcome{}, fmt.Errorf("read external pack output: %w", scanErr)
	}

	if !completed {
		return Outcome{}, fmt.Errorf("external pack completed without a completion record")
	}

	outcome.Diagnostics = diagnostics
	return outcome, nil
}

func decodeDiagnostic(record envelope, line int) (diagnostic style.Diagnostic, err error) {
	if strings.TrimSpace(record.Message) == "" {
		return style.Diagnostic{}, fmt.Errorf(
			"external pack diagnostic at line %d is missing a message", line)
	}

	rng := style.Range{
		Start: position(record.Start),
		End:   position(record.End),
	}

	if err = style.VerifyRange(record.File, rng); err != nil {
		return style.Diagnostic{}, fmt.Errorf("external pack diagnostic at line %d: %w", line, err)
	}

	return style.Diagnostic{
		Code:    record.Code,
		Message: record.Message,
		File:    record.File,
		Range:   rng,
		HelpURL: record.HelpURL,
	}, nil
}

func position(reference *pos) (spot style.Position) {
	if reference == nil {
		return style.Position{}
	}
	return style.Position{Line: reference.Line, Column: reference.Column}
}
