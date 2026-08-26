package external_test

import (
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/pack/external"
	"github.com/wbd2023/quill/internal/style"
)

/* -------------------------------------- Response Parsing -------------------------------------- */

func TestParseResponseAcceptsDiagnosticsAndCompletion(t *testing.T) {
	t.Parallel()

	stdout := strings.Join([]string{
		`{"type":"diagnostic","code":"db-access",` +
			`"message":"Use the repository adapter.",` +
			`"file":"internal/service/users.go",` +
			`"start":{"line":42,"column":5},` +
			`"end":{"line":42,"column":18}}`,
		`{"type":"diagnostic","code":"db-access",` +
			`"message":"Second finding.",` +
			`"file":"internal/service/orders.go",` +
			`"start":{"line":7,"column":1}}`,
		`{"type":"complete","success":true}`,
		"",
	}, "\n")

	outcome, err := external.ParseResponse(stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !outcome.Succeeded {
		t.Fatalf("expected success completion")
	}
	if len(outcome.Diagnostics) != 2 {
		t.Fatalf("expected 2 diagnostics, got %d", len(outcome.Diagnostics))
	}

	first := outcome.Diagnostics[0]
	if first.Code != "db-access" || first.File != "internal/service/users.go" {
		t.Fatalf("unexpected first diagnostic: %+v", first)
	}
	if first.Range.Start.Line != 42 || first.Range.Start.Column != 5 {
		t.Fatalf("unexpected start: %+v", first.Range.Start)
	}
	if first.Range.End.Line != 42 || first.Range.End.Column != 18 {
		t.Fatalf("unexpected end: %+v", first.Range.End)
	}
}

func TestParseResponsePreservesHelpURL(t *testing.T) {
	t.Parallel()

	const helpURL = "https://example.invalid/rules/db-access"
	stdout := strings.Join([]string{
		`{"type":"diagnostic","code":"db-access",` +
			`"message":"Use the repository adapter.",` +
			`"file":"internal/service/users.go",` +
			`"start":{"line":42,"column":5},` +
			`"end":{"line":42,"column":18},` +
			`"help_url":"` + helpURL + `"}`,
		`{"type":"complete","success":true}`,
		"",
	}, "\n")

	outcome, err := external.ParseResponse(stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(outcome.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(outcome.Diagnostics))
	}
	if got := outcome.Diagnostics[0].HelpURL; got != helpURL {
		t.Fatalf("expected help_url %q preserved on the diagnostic, got %q", helpURL, got)
	}
}

func TestParseResponseAcceptsRepositoryLevelDiagnostic(t *testing.T) {
	t.Parallel()

	stdout := `{"type":"diagnostic","message":"Pack-wide finding."}` + "\n" +
		`{"type":"complete","success":true}` + "\n"

	outcome, err := external.ParseResponse(stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(outcome.Diagnostics) != 1 {
		t.Fatalf("expected 1 diagnostic, got %d", len(outcome.Diagnostics))
	}
	if outcome.Diagnostics[0].File != "" || !outcome.Diagnostics[0].Range.IsUnknown() {
		t.Fatalf("expected repository-level diagnostic, got %+v", outcome.Diagnostics[0])
	}
}

func TestParseResponseBlankDiagnosticMessageRejected(t *testing.T) {
	t.Parallel()

	stdout := `{"type":"diagnostic","message":"   "}` + "\n" +
		`{"type":"complete","success":true}` + "\n"

	if _, err := external.ParseResponse(stdout); err == nil {
		t.Fatal("expected error for blank diagnostic message")
	}
}

func TestParseResponseCompletionFailure(t *testing.T) {
	t.Parallel()

	stdout := `{"type":"complete","success":false,"error":"config field missing"}` + "\n"

	outcome, err := external.ParseResponse(stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if outcome.Succeeded {
		t.Fatal("expected failure completion")
	}
	if outcome.Error != "config field missing" {
		t.Fatalf("unexpected error message: %q", outcome.Error)
	}
}

func TestParseResponseProtocolFailures(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		stdout string
	}{
		{"missing completion", `{"type":"diagnostic","message":"x"}` + "\n"},
		{"empty output", ""},
		{
			"duplicate completion",
			`{"type":"complete","success":true}` + "\n" +
				`{"type":"complete","success":true}` + "\n",
		},
		{
			"diagnostic after completion",
			`{"type":"complete","success":true}` + "\n" +
				`{"type":"diagnostic","message":"late"}` + "\n",
		},
		{"malformed json", `not json at all` + "\n"},
		{
			"unknown record type",
			`{"type":"progress","message":"halfway"}` + "\n" +
				`{"type":"complete","success":true}` + "\n",
		},
		{
			"absolute path rejected",
			`{"type":"diagnostic","message":"x","file":"/etc/passwd"}` + "\n" +
				`{"type":"complete","success":true}` + "\n",
		},
		{
			"path escape rejected",
			`{"type":"diagnostic","message":"x","file":"../secret"}` + "\n" +
				`{"type":"complete","success":true}` + "\n",
		},
		{
			"inverted range rejected",
			`{"type":"diagnostic","message":"x","file":"a.go",` +
				`"start":{"line":10,"column":5},` +
				`"end":{"line":5,"column":2}}` + "\n" +
				`{"type":"complete","success":true}` + "\n",
		},
		{
			"column without line rejected",
			`{"type":"diagnostic","message":"x","file":"a.go",` +
				`"start":{"line":0,"column":5}}` + "\n" +
				`{"type":"complete","success":true}` + "\n",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if _, err := external.ParseResponse(testCase.stdout); err == nil {
				t.Fatalf("expected error for %s", testCase.name)
			}
		})
	}
}

func TestParseResponseSkipsBlankLines(t *testing.T) {
	t.Parallel()

	stdout := "\n\n" + `{"type":"complete","success":true}` + "\n\n"

	outcome, err := external.ParseResponse(stdout)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !outcome.Succeeded {
		t.Fatal("expected success")
	}
	if len(outcome.Diagnostics) != 0 {
		t.Fatalf("expected 0 diagnostics, got %d", len(outcome.Diagnostics))
	}
}

/* -------------------------------------- Request Encoding -------------------------------------- */

func TestEncodeRequestRoundTripsProtocol(t *testing.T) {
	t.Parallel()

	payload, err := external.EncodeRequest(external.Request{
		Protocol:  external.ProtocolVersion,
		Operation: "check",
		PackID:    "company",
		RuleID:    "company/rule",
		CheckID:   "rule",
		Files:     []string{"a.go", "b.go"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(string(payload), `"protocol":"quill-pack-v1"`) {
		t.Fatalf("payload missing protocol: %s", payload)
	}
	if !strings.Contains(string(payload), `"operation":"check"`) {
		t.Fatalf("payload missing operation: %s", payload)
	}
}

/* ------------------------------------- Range Verification ------------------------------------- */

func TestVerifyRangeContractForExternalDiagnostics(t *testing.T) {
	t.Parallel()

	// Guards that the external protocol ingestion path honours the style boundary verifier: a
	// diagnostic whose range would be admitted without verification is a trust-boundary violation.
	if err := style.VerifyRange("../escape", style.Range{}); err == nil {
		t.Fatal("VerifyRange must reject path escapes")
	}
}
