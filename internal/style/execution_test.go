package style_test

import (
	"reflect"
	"slices"
	"testing"

	"github.com/wbd2023/quill/internal/style"
)

/* ------------------------------------------ Describe ------------------------------------------ */

func TestDescribeReturnsCorrectRequirements(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		template style.Template
		want     style.TemplateRequirements
	}{
		{
			name:     "toolchain",
			template: style.ToolchainCheck{ToolIDs: []string{"go"}},
			want:     style.TemplateRequirements{ToolIDs: []string{"go"}},
		},
		{
			name:     "profile",
			template: style.ProfileCheck{Check: "config"},
			want:     style.TemplateRequirements{},
		},
		{
			name:     "file command",
			template: style.FileCommand{ToolID: "tool", FileSet: "go"},
			want:     style.TemplateRequirements{ToolIDs: []string{"tool"}, FileSet: "go"},
		},
		{
			name:     "repository scan",
			template: style.RepositoryScan{FileSet: "text"},
			want:     style.TemplateRequirements{FileSet: "text"},
		},
		{
			name:     "external check",
			template: style.ExternalCheck{FileSet: "ext"},
			want:     style.TemplateRequirements{FileSet: "ext"},
		},
		{
			name:     "target command",
			template: style.TargetCommandTemplate{ToolIDs: []string{"go"}, Language: "go"},
			want: style.TemplateRequirements{
				ToolIDs:        []string{"go"},
				NeedsTargets:   true,
				TargetLanguage: "go",
			},
		},
		{
			name:     "target check",
			template: style.TargetCheckTemplate{ToolIDs: []string{"go"}, Language: "go"},
			want: style.TemplateRequirements{
				ToolIDs:         []string{"go"},
				NeedsTargets:    true,
				TargetLanguage:  "go",
				NeedsCheckPaths: true,
			},
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

// TestDescribeDoesNotAliasOwnedSlices proves Describe defensively copies mutable slices so a
// caller cannot mutate the source Template through the returned TemplateRequirements.
func TestDescribeDoesNotAliasOwnedSlices(t *testing.T) {
	t.Parallel()

	template := style.ToolchainCheck{ToolIDs: []string{"go", "gofmt"}}
	requirements := style.Describe(template)

	requirements.ToolIDs[0] = "mutated"
	if got := template.ToolIDs[0]; got != "go" {
		t.Fatalf("source ToolIDs mutated via Describe: %q", got)
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
		job := template.Bind([]string{"target1", "target2"})

		if job.Action != "format" {
			t.Fatalf("Action = %q, want %q", job.Action, "format")
		}
		if !slices.Equal(job.Targets, []string{"target1", "target2"}) {
			t.Fatalf("Targets = %v, want %v", job.Targets, []string{"target1", "target2"})
		}
		if !slices.Equal(job.ToolIDs, []string{"go"}) {
			t.Fatalf("ToolIDs = %v, want %v", job.ToolIDs, []string{"go"})
		}
	})

	t.Run("target check binds targets into job", func(t *testing.T) {
		t.Parallel()

		template := style.TargetCheckTemplate{
			ToolIDs:  []string{"go"},
			Check:    "style",
			Language: "go",
		}
		job := template.Bind([]string{"target1"})

		if job.Check != "style" {
			t.Fatalf("Check = %q, want %q", job.Check, "style")
		}
		if !slices.Equal(job.Targets, []string{"target1"}) {
			t.Fatalf("Targets = %v, want %v", job.Targets, []string{"target1"})
		}
	})

	t.Run("bind copies targets independently of the resolver slice", func(t *testing.T) {
		t.Parallel()

		template := style.TargetCommandTemplate{ToolIDs: []string{"go"}, Language: "go"}
		resolver := []string{"target1", "target2"}
		job := template.Bind(resolver)

		resolver[0] = "mutated"
		if got := job.Targets[0]; got != "target1" {
			t.Fatalf("job Targets mutated via resolver slice: %q", got)
		}
	})

	t.Run("bind copies tool ids independently of the template slice", func(t *testing.T) {
		t.Parallel()

		template := style.TargetCommandTemplate{ToolIDs: []string{"go"}, Language: "go"}
		job := template.Bind([]string{"target1"})

		template.ToolIDs[0] = "mutated"
		if got := job.ToolIDs[0]; got != "go" {
			t.Fatalf("job ToolIDs mutated via template slice: %q", got)
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
			job:  style.ToolchainCheck{ToolIDs: []string{"go", "gofmt"}},
			want: []string{"go", "gofmt"},
		},
		{
			name: "file command",
			job:  style.FileCommand{ToolID: "grep"},
			want: []string{"grep"},
		},
		{
			name: "file command empty tool id",
			job:  style.FileCommand{},
			want: nil,
		},
		{
			name: "target command job",
			job:  style.TargetCommandJob{ToolIDs: []string{"go"}},
			want: []string{"go"},
		},
		{
			name: "target check job",
			job:  style.TargetCheckJob{ToolIDs: []string{"go"}},
			want: []string{"go"},
		},
		{
			name: "profile",
			job:  style.ProfileCheck{},
			want: nil,
		},
		{
			name: "repository scan",
			job:  style.RepositoryScan{},
			want: nil,
		},
		{
			name: "external check",
			job:  style.ExternalCheck{},
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

// TestToolIDsDoesNotAliasOwnedSlices proves ToolIDs defensively copies so a caller cannot mutate
// the source Job through the returned slice.
func TestToolIDsDoesNotAliasOwnedSlices(t *testing.T) {
	t.Parallel()

	job := style.ToolchainCheck{ToolIDs: []string{"go", "gofmt"}}
	ids := style.ToolIDs(job)

	ids[0] = "mutated"
	if got := job.ToolIDs[0]; got != "go" {
		t.Fatalf("source ToolIDs mutated via ToolIDs: %q", got)
	}
}

/* -------------------------------------- Phase Membership -------------------------------------- */

// TestNonTargetValuesSatisfyBothInterfaces proves the five non-target families implement both
// Template and Job because Profile binding adds nothing to them.
func TestNonTargetValuesSatisfyBothInterfaces(t *testing.T) {
	t.Parallel()

	values := []style.Job{
		style.ToolchainCheck{ToolIDs: []string{"go"}},
		style.ProfileCheck{Check: "config"},
		style.FileCommand{ToolID: "tool", FileSet: "go"},
		style.RepositoryScan{Scanner: "secrets"},
		style.ExternalCheck{CheckID: "ext"},
	}

	for _, value := range values {
		if _, ok := value.(style.Template); !ok {
			t.Fatalf("%T satisfies Job but not Template", value)
		}
	}
}

// TestTargetTemplatesAreNotJobs proves a target Template cannot reach a driver before binding: it
// must be converted to its paired target Job first.
func TestTargetTemplatesAreNotJobs(t *testing.T) {
	t.Parallel()

	templates := []style.Template{
		style.TargetCommandTemplate{ToolIDs: []string{"go"}, Language: "go"},
		style.TargetCheckTemplate{ToolIDs: []string{"go"}, Language: "go"},
	}

	for _, template := range templates {
		if _, ok := template.(style.Job); ok {
			t.Fatalf("%T must not satisfy Job before binding", template)
		}
	}
}

// TestTargetJobsAreNotTemplates proves a bound target Job cannot be mistaken for an unbound
// Template, so metadata and compilation never treat a compiled value as a declaration.
func TestTargetJobsAreNotTemplates(t *testing.T) {
	t.Parallel()

	jobs := []style.Job{
		style.TargetCommandJob{ToolIDs: []string{"go"}, Language: "go", Targets: []string{"t"}},
		style.TargetCheckJob{ToolIDs: []string{"go"}, Language: "go", Targets: []string{"t"}},
	}

	for _, job := range jobs {
		if _, ok := job.(style.Template); ok {
			t.Fatalf("%T must not satisfy Template", job)
		}
	}
}

// TestTemplatesAreSealed proves no foreign type can satisfy Template, since describe is
// unexported. A reflective zero-value of an invented type must fail the assertion.
func TestTemplatesAreSealed(t *testing.T) {
	t.Parallel()

	type foreign struct{}
	var _ style.Job = style.ToolchainCheck{} // compile-time confirmation the real type qualifies

	var impostor foreign
	if _, ok := any(impostor).(style.Template); ok {
		t.Fatalf("%T must not satisfy the sealed Template interface", impostor)
	}
	if _, ok := any(impostor).(style.Job); ok {
		t.Fatalf("%T must not satisfy the sealed Job interface", impostor)
	}

	// Sanity: reflect cannot find an unexported describe/toolIDs method on the foreign type.
	if method, found := reflect.TypeOf(impostor).MethodByName("describe"); found {
		t.Fatalf("foreign type unexpectedly exposes describe: %v", method)
	}
}
