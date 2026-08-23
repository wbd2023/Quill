package drivers

import (
	"context"
	"testing"

	"github.com/wbd2023/quill/internal/execution"
	"github.com/wbd2023/quill/internal/profile"
	"github.com/wbd2023/quill/internal/style"
	"github.com/wbd2023/quill/internal/toolchain"
)

func TestToolchainDriverReportsMissingStatus(t *testing.T) {
	run := execution.NewRunContext(
		t.TempDir(),
		style.Scope("all"),
		profile.Profile{},
		style.Plan{},
		nil,
		nil,
		nil,
	)

	result, err := ToolchainDriver(
		context.Background(),
		run,
		style.Rule{},
		style.ToolchainCheck{ToolIDs: []string{"go"}},
		toolchain.StatusMap{},
	)
	if err != nil {
		t.Fatalf("ToolchainDriver: %v", err)
	}
	if len(result.Diagnostics) != 1 {
		t.Fatalf("expected one diagnostic, got %#v", result.Diagnostics)
	}
	if got := result.Diagnostics[0]; got.Code != "toolchain/invalid" ||
		got.Message != "go: no inspection status" {
		t.Fatalf("unexpected diagnostic: %#v", got)
	}
}
