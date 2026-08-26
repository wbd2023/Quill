package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/wbd2023/quill/internal/report"
	"github.com/wbd2023/quill/internal/testutil"
)

/* ---------------------------------------- Machine Mode ---------------------------------------- */

// decodedEnvelope captures the protocol envelope fields shared by every machine response. Result
// is kept raw so command-specific payloads can be inspected independently.
type decodedEnvelope struct {
	SchemaVersion int             `json:"schema_version"`
	QuillVersion  string          `json:"quill_version"`
	Command       string          `json:"command"`
	Status        string          `json:"status"`
	Result        json.RawMessage `json:"result,omitempty"`
	Error         *decodedError   `json:"error,omitempty"`
}

type decodedError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

/* ------------------------------------- Response Envelopes ------------------------------------- */

func TestMachineCoverageEmitsDecodedEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), machineCoverageArgs(t))
	if exitCode != 0 {
		t.Fatalf("expected exit 0 for coverage, got %d", exitCode)
	}

	envelope := assertResultEnvelope(t, stdout.Bytes(), "coverage")

	var payload struct {
		Report struct {
			Requirements []struct {
				ID string `json:"id"`
			} `json:"requirements"`
		} `json:"report"`
	}
	if err := json.Unmarshal(envelope.Result, &payload); err != nil {
		t.Fatalf("decode coverage result: %v", err)
	}

	// The repository STYLE.md is the source of truth; coverage must return its requirements.
	if len(payload.Report.Requirements) == 0 {
		t.Fatalf("coverage result carried no requirements: %s", envelope.Result)
	}
}

func TestMachineDoctorEmitsDecodedEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), []string{
		"doctor", "--format", "json", "--repository-root", testutil.RepositoryRoot(t),
	})

	// doctor completes regardless of tool validity; an invalid toolchain still reports status ok
	// with a nonzero exit status per the documented "status ok + nonzero exit" rule.
	if exitCode != 0 && exitCode != 1 {
		t.Fatalf("expected exit 0 or 1 for doctor, got %d", exitCode)
	}

	envelope := assertResultEnvelope(t, stdout.Bytes(), "doctor")

	var payload struct {
		AllValid bool `json:"all_valid"`
	}
	if err := json.Unmarshal(envelope.Result, &payload); err != nil {
		t.Fatalf("decode doctor result: %v", err)
	}
}

func TestMachineCheckInvalidArgumentEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), []string{
		"check", "--format", "json", "--mode", "invalid",
		"--repository-root", testutil.RepositoryRoot(t),
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code %d for invalid argument, got %d",
			usageExitCode, exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "check", report.ErrorCodeInvalidArgument)
}

func TestMachineCheckUsesFormatBeforePositionalArgument(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), []string{
		"check", "--format", "json", "unexpected", "--format", "text",
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code %d for positional argument, got %d",
			usageExitCode, exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "check", report.ErrorCodeInvalidArgument)
}

func TestMachineFixInvalidArgumentEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), []string{
		"fix", "--format", "json", "unexpected", "--repository-root", testutil.RepositoryRoot(t),
	})
	if exitCode != usageExitCode {
		t.Fatalf("expected usage exit code %d for positional argument, got %d",
			usageExitCode, exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "fix", report.ErrorCodeInvalidArgument)
}

func TestMachineCoverageOperationFailedEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), []string{
		"coverage", "--format", "json", "--repository-root", t.TempDir(),
	})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for operation failure, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "coverage", report.ErrorCodeOperationFailed)
}

func TestMachineAutoDetectedMissingRootIsOperationFailed(t *testing.T) {
	t.Chdir(t.TempDir())

	tool, stdout, _ := newMachineCLI()
	exitCode := tool.Run(context.Background(), []string{"check", "--format", "json"})
	if exitCode != 1 {
		t.Fatalf("expected operation exit 1, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "check", report.ErrorCodeOperationFailed)
}

func TestMachineInstallOperationFailedEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), []string{
		"install", "--format", "json", "--repository-root", t.TempDir(),
	})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for operation failure, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "install", report.ErrorCodeOperationFailed)
}

func TestMachineLockOperationFailedEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.Run(context.Background(), []string{
		"lock", "--format", "json", "--repository-root", t.TempDir(),
	})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for operation failure, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "lock", report.ErrorCodeOperationFailed)
}

func TestMachineCommandTimeoutEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.reportCommandError(
		"check",
		report.FormatJSON,
		context.DeadlineExceeded,
	)
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for command timeout, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "check", report.ErrorCodeOperationFailed)
}

func TestMachineCommandErrorUsesReturnedCancellation(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	exitCode := tool.reportCommandError(
		"check",
		report.FormatJSON,
		fmt.Errorf("operation: %w", context.Canceled),
	)
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for cancellation, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "check", report.ErrorCodeCancelled)
}

/* ---------------------------------------- Cancellation ---------------------------------------- */

func TestMachineCoverageCancelledEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	exitCode := tool.Run(ctx, machineCoverageArgs(t))
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for cancelled operation, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "coverage", report.ErrorCodeCancelled)
}

func TestMachineCheckCancelledEnvelope(t *testing.T) {
	tool, stdout, _ := newMachineCLI()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	exitCode := tool.Run(ctx, []string{
		"check", "--format", "json", "--repository-root", testutil.RepositoryRoot(t),
	})
	if exitCode != 1 {
		t.Fatalf("expected exit 1 for cancelled operation, got %d", exitCode)
	}

	assertErrorEnvelope(t, stdout.Bytes(), "check", report.ErrorCodeCancelled)
}

/* ---------------------------------------- Stdout Purity --------------------------------------- */

func TestMachineStdoutCarriesOnlyEnvelope(t *testing.T) {
	cases := []struct {
		name string
		args []string
	}{
		{name: "coverage success", args: machineCoverageArgs(t)},
		{
			name: "check invalid argument",
			args: []string{"check", "--format", "json", "--mode", "invalid",
				"--repository-root", testutil.RepositoryRoot(t)},
		},
		{
			name: "coverage operation failed",
			args: []string{"coverage", "--format", "json", "--repository-root", t.TempDir()},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			tool, stdout, stderr := newMachineCLI()

			_ = tool.Run(context.Background(), test.args)

			output := stdout.Bytes()
			// json.Valid guarantees the entire stdout is exactly one JSON document (it rejects
			// trailing content), so diagnostics, usage, and progress must not share the stream.
			if !json.Valid(output) {
				t.Fatalf("stdout must be exactly one JSON document, got:\n%s", output)
			}

			if strings.Contains(stderr.String(), `"schema_version"`) {
				t.Fatalf("envelope must not leak to stderr:\n%s", stderr.String())
			}
		})
	}
}

/* ---------------------------------------- Test Harness ---------------------------------------- */

func newMachineCLI() (runner Runner, stdout *bytes.Buffer, stderr *bytes.Buffer) {
	stdout = &bytes.Buffer{}
	stderr = &bytes.Buffer{}
	runner = New(stdout, stderr, "test-version")
	return runner, stdout, stderr
}

/* ------------------------------------------- Helpers ------------------------------------------ */

func decodeEnvelope(t *testing.T, output []byte) (envelope decodedEnvelope) {
	t.Helper()

	if !json.Valid(output) {
		t.Fatalf("stdout must be exactly one JSON document, got:\n%s", output)
	}

	if err := json.Unmarshal(output, &envelope); err != nil {
		t.Fatalf("decode envelope: %v\n%s", err, output)
	}

	return envelope
}

func assertResultEnvelope(
	t *testing.T,
	output []byte,
	command string,
) (envelope decodedEnvelope) {
	t.Helper()

	envelope = decodeEnvelope(t, output)
	if envelope.SchemaVersion != report.SchemaVersion {
		t.Fatalf("schema_version = %d, want %d", envelope.SchemaVersion, report.SchemaVersion)
	}

	if envelope.QuillVersion != "test-version" {
		t.Fatalf("quill_version = %q, want %q", envelope.QuillVersion, "test-version")
	}

	if envelope.Command != command {
		t.Fatalf("command = %q, want %q", envelope.Command, command)
	}

	if envelope.Status != report.StatusOK {
		t.Fatalf("status = %q, want %q", envelope.Status, report.StatusOK)
	}

	if envelope.Error != nil {
		t.Fatalf("success envelope must not carry an error: %+v", envelope.Error)
	}

	if len(envelope.Result) == 0 || string(envelope.Result) == "null" {
		t.Fatalf("success envelope must carry a result payload")
	}

	return envelope
}

func assertErrorEnvelope(
	t *testing.T,
	output []byte,
	command string,
	code string,
) (envelope decodedEnvelope) {
	t.Helper()

	envelope = decodeEnvelope(t, output)
	if envelope.SchemaVersion != report.SchemaVersion {
		t.Fatalf("schema_version = %d, want %d", envelope.SchemaVersion, report.SchemaVersion)
	}

	if envelope.QuillVersion != "test-version" {
		t.Fatalf("quill_version = %q, want %q", envelope.QuillVersion, "test-version")
	}

	if envelope.Command != command {
		t.Fatalf("command = %q, want %q", envelope.Command, command)
	}

	if envelope.Status != report.StatusError {
		t.Fatalf("status = %q, want %q", envelope.Status, report.StatusError)
	}

	if envelope.Result != nil {
		t.Fatalf("error envelope must not carry a result, got %s", envelope.Result)
	}

	if envelope.Error == nil {
		t.Fatal("error envelope missing error payload")
	}

	if envelope.Error.Code != code {
		t.Fatalf("error code = %q, want %q", envelope.Error.Code, code)
	}

	if strings.TrimSpace(envelope.Error.Message) == "" {
		t.Fatal("error envelope missing message")
	}

	return envelope
}

func machineCoverageArgs(t *testing.T) (arguments []string) {
	t.Helper()
	return []string{"coverage", "--format", "json", "--repository-root", testutil.RepositoryRoot(t)}
}
