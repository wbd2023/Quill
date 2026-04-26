package profile

import "testing"

/* ------------------------------------------- Loading ------------------------------------------ */

func TestParseProfileRejectsUnknownKeys(t *testing.T) {
	t.Parallel()

	_, err := parseProfile(`profile_version = 1
unknown_key = true
`)
	if err == nil {
		t.Fatal("expected unknown key to fail")
	}
}

func TestParseProfileDecodesMultilineArrays(t *testing.T) {
	t.Parallel()

	config, err := parseProfile(`profile_version = 1

[paths]
go_source = [
	"cmd/",
	"internal/",
]
`)
	if err != nil {
		t.Fatalf("parseProfile: %v", err)
	}

	got := config.Paths.Patterns("go_source")
	if len(got) != 2 || got[0] != "cmd/" || got[1] != "internal/" {
		t.Fatalf("paths.go_source = %v", got)
	}
}

func TestParseProfileDecodesSectionHeaderPolicy(t *testing.T) {
	t.Parallel()

	config, err := parseProfile(`profile_version = 1

[formatting.section_headers]
required_min_lines = 100
short_file_max_lines = 79
overuse_header_count = 6
generic_names = ["Check", "Checks"]
structural_names = ["Types", "Helpers"]
`)
	if err != nil {
		t.Fatalf("parseProfile: %v", err)
	}

	headers := config.Formatting.SectionHeaders
	if headers.RequiredMinLines != 100 ||
		headers.ShortFileMaxLines != 79 ||
		headers.OveruseCount != 6 ||
		len(headers.GenericNames) != 2 ||
		headers.GenericNames[0] != "Check" ||
		len(headers.StructuralNames) != 2 {
		t.Fatalf("unexpected section header policy: %#v", headers)
	}
}

func TestParseProfileRejectsOldScanRootFields(t *testing.T) {
	t.Parallel()

	_, err := parseProfile(`profile_version = 1

[repository]
app_scan_roots = ["cmd"]
`)
	if err == nil {
		t.Fatal("expected old repository scan-root field to fail")
	}
}

func TestLoadRejectsMissingConfiguredRootMarker(t *testing.T) {
	repoRoot := t.TempDir()
	writeFile(t, repoRoot, "STYLE.md", "# Test Style Guide\n")
	writeFile(
		t,
		repoRoot,
		"style.toml",
		`profile_version = 1

[rule_packs]
enabled = ["repository"]

[repository]
root_markers = ["STYLE.md", "style.toml", "PROJECT.marker"]
default_scope = "all"
global_exclusions = [".git"]
generated_marker = "DO NOT EDIT."
generated_probe_limit = 128

[repository.scopes]
app = ["cmd"]
tools = ["tools"]
all = ["."]

[styleguide]
path = "STYLE.md"
requirement_id_format = "section_slug"

[formatting.section_headers]
required_min_lines = 100
short_file_max_lines = 79
overuse_header_count = 6
generic_names = ["Check", "Checks", "Misc", "Other"]
structural_names = ["Types", "Constants", "Helpers"]

[imports]
local_prefix = "example.com/test"

[paths]
go_source = ["cmd/"]
application_port = ["internal/app/ports/"]
concrete_infra = ["internal/adapters/"]
domain = ["internal/domain/"]
domain_errors = ["internal/domain/errors.go"]
test_mocks = ["internal/testsupport/mocks/"]

[[file_sets]]
name = "markdown"
extensions = [".md"]
files = { app = ["STYLE.md"] }
prefixes = { tools = ["tools/"] }

[[language.backends]]
name = "application_go"
language = "go"
scope = "app"
workdir = "."
format_paths = ["cmd"]
style_paths = ["cmd"]

[naming]
go_type_suffix_forbidden = ["Repository"]
go_type_suffix_preferred = "Store"
go_identifier_suffix_forbidden = ["Repository"]
go_identifier_suffix_preferred = "Store"
shell_forbidden_assignments = ["NC"]
shell_preferred_assignment = "COLOUR_RESET"

[control_plane]
quality_file = "mk/quality.mk"

[[control_plane.variable_contracts]]
name = "LINT_REQUIRED_ARGS"
value = "--mode required"

[[control_plane.target_contracts]]
name = "lint"
recipe_line = "@$(STYLE_CMD) check $(LINT_FULL_ARGS)"

[[architecture.layers]]
name = "domain"
package_roots = ["internal/domain"]
may_import = ["domain"]

[[rules]]
rule_id = "naming/vocabulary"
level = "required"
scope = "all"
requirement_ids = ["3.3.use-repository-name"]
`,
	)

	if _, err := Load(repoRoot); err == nil {
		t.Fatal("expected missing configured root marker to fail")
	}
}
