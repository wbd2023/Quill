package style_test

import (
	"slices"
	"testing"

	"github.com/wbd2023/Quill/internal/style"
)

func TestDescribeReturnsCorrectRequirements(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		template style.Template
		want     style.Requirements
	}{
		{
			name:     "toolchain",
			template: style.ToolchainExecution{ToolIDs: []string{"go"}},
			want:     style.Requirements{ToolIDs: []string{"go"}},
		},
		{
			name:     "profile",
			template: style.ProfileExecution{Check: "config"},
			want:     style.Requirements{},
		},
		{
			name:     "file command",
			template: style.FileCommandExecution{ToolID: "tool", FileSet: "go"},
			want:     style.Requirements{ToolIDs: []string{"tool"}, FileSet: "go"},
		},
		{
			name:     "target command",
			template: style.TargetCommandTemplate{ToolIDs: []string{"go"}, Language: "go"},
			want: style.Requirements{
				ToolIDs:        []string{"go"},
				NeedsTargets:   true,
				TargetLanguage: "go",
			},
		},
		{
			name:     "target check",
			template: style.TargetCheckTemplate{ToolIDs: []string{"go"}, Language: "go"},
			want: style.Requirements{
				ToolIDs:         []string{"go"},
				NeedsTargets:    true,
				TargetLanguage:  "go",
				NeedsCheckPaths: true,
			},
		},
		{
			name:     "repository scan",
			template: style.RepositoryScanExecution{FileSet: "text"},
			want:     style.Requirements{FileSet: "text"},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := style.Describe(test.template)
			if !slices.Equal(got.ToolIDs, test.want.ToolIDs) {
				t.Fatalf("ToolIDs = %v, want %v", got.ToolIDs, test.want.ToolIDs)
			}
			if got.FileSet != test.want.FileSet {
				t.Fatalf("FileSet = %q, want %q", got.FileSet, test.want.FileSet)
			}
			if got.NeedsTargets != test.want.NeedsTargets {
				t.Fatalf("NeedsTargets = %t, want %t", got.NeedsTargets, test.want.NeedsTargets)
			}
			if got.TargetLanguage != test.want.TargetLanguage {
				t.Fatalf("TargetLanguage = %q, want %q",
					got.TargetLanguage, test.want.TargetLanguage)
			}
			if got.NeedsCheckPaths != test.want.NeedsCheckPaths {
				t.Fatalf("NeedsCheckPaths = %t, want %t",
					got.NeedsCheckPaths, test.want.NeedsCheckPaths)
			}
		})
	}
}

/* -------------------------------------------- Bind -------------------------------------------- */

func TestBindProducesCorrectJobs(t *testing.T) {
	t.Parallel()

	t.Run("target command binds targets into job", func(t *testing.T) {
		t.Parallel()

		template := style.TargetCommandTemplate{
			ToolIDs:  []string{"go"},
			Action:   "format",
			Language: "go",
		}
		job := style.Bind(template, []string{"target1", "target2"})

		concrete, ok := job.(style.TargetCommandJob)
		if !ok {
			t.Fatalf("expected TargetCommandJob, got %T", job)
		}
		if concrete.Action != "format" {
			t.Fatalf("Action = %q, want %q", concrete.Action, "format")
		}
		if !slices.Equal(concrete.Targets, []string{"target1", "target2"}) {
			t.Fatalf("Targets = %v, want %v", concrete.Targets, []string{"target1", "target2"})
		}
	})

	t.Run("target check binds targets into job", func(t *testing.T) {
		t.Parallel()

		template := style.TargetCheckTemplate{
			ToolIDs:  []string{"go"},
			Check:    "style",
			Language: "go",
		}
		job := style.Bind(template, []string{"target1"})

		concrete, ok := job.(style.TargetCheckJob)
		if !ok {
			t.Fatalf("expected TargetCheckJob, got %T", job)
		}
		if concrete.Check != "style" {
			t.Fatalf("Check = %q, want %q", concrete.Check, "style")
		}
		if !slices.Equal(concrete.Targets, []string{"target1"}) {
			t.Fatalf("Targets = %v, want %v", concrete.Targets, []string{"target1"})
		}
	})

	t.Run("non-target template returns self as job", func(t *testing.T) {
		t.Parallel()

		template := style.ToolchainExecution{ToolIDs: []string{"go"}}
		job := style.Bind(template, nil)

		if _, ok := job.(style.ToolchainExecution); !ok {
			t.Fatalf("expected ToolchainExecution, got %T", job)
		}
	})
}

/* ------------------------------------------ Tool IDs ------------------------------------------ */

func TestToolIDsReturnsCorrectIDs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name string
		job  style.Job
		want []string
	}{
		{
			name: "toolchain",
			job:  style.ToolchainExecution{ToolIDs: []string{"go", "gofmt"}},
			want: []string{"go", "gofmt"},
		},
		{
			name: "file command",
			job:  style.FileCommandExecution{ToolID: "grep"},
			want: []string{"grep"},
		},
		{
			name: "file command empty tool id",
			job:  style.FileCommandExecution{},
			want: nil,
		},
		{
			name: "target command job",
			job:  style.TargetCommandJob{ToolIDs: []string{"go"}},
			want: []string{"go"},
		},
		{
			name: "profile",
			job:  style.ProfileExecution{},
			want: nil,
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got := style.ToolIDs(test.job)
			if !slices.Equal(got, test.want) {
				t.Fatalf("ToolIDs = %v, want %v", got, test.want)
			}
		})
	}
}
