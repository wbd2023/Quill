package report

import (
	"bytes"
	"encoding/json"
	"testing"

	"ciphera/tools/internal/contract"
	"ciphera/tools/internal/runtime"
)

/* -------------------------------------- Toolchain Output -------------------------------------- */

func TestWriteToolchainText(t *testing.T) {
	var buffer bytes.Buffer

	allValid, err := WriteToolchain(&buffer, FormatText, NewToolchainView(ToolchainResult{
		Statuses: []runtime.ToolStatus{
			{
				Tool:    contract.Tool{Name: "Go"},
				Version: "1.24.5",
				Valid:   true,
			},
			{
				Tool:    contract.Tool{Name: "markdownlint"},
				Version: "0.48.0",
				Valid:   false,
				Issue:   "requires pinned version 0.45.0",
			},
		},
	}))
	if err != nil {
		t.Fatalf("WriteToolchain: %v", err)
	}

	if allValid {
		t.Fatal("expected invalid toolchain")
	}

	output := buffer.String()
	if output != readGoldenOutput(t, "toolchain.txt") {
		t.Fatalf("unexpected toolchain output:\n%s", output)
	}
}

func TestWriteToolchainJSON(t *testing.T) {
	var buffer bytes.Buffer

	view := NewToolchainView(ToolchainResult{
		Statuses: []runtime.ToolStatus{
			{
				Tool:  contract.Tool{Name: "Go"},
				Valid: true,
			},
			{
				Tool:  contract.Tool{Name: "markdownlint"},
				Valid: false,
				Issue: "requires pinned version 0.45.0",
			},
		},
	})
	allValid, err := WriteToolchain(&buffer, FormatJSON, view)
	if err != nil {
		t.Fatalf("WriteToolchain: %v", err)
	}

	if allValid {
		t.Fatal("expected invalid toolchain")
	}

	var envelope struct {
		Toolchain ToolchainView `json:"toolchain"`
	}
	if err := json.Unmarshal(buffer.Bytes(), &envelope); err != nil {
		t.Fatalf("decode toolchain json: %v", err)
	}

	if envelope.Toolchain.AllValid {
		t.Fatal("expected all_valid=false in JSON output")
	}
}
