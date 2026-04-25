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

	policy, err := parseProfile(`profile_version = 1

[paths]
app = [
	"cmd/",
	"internal/",
]
`)
	if err != nil {
		t.Fatalf("parseProfile: %v", err)
	}

	got := policy.Paths.Patterns("app")
	if len(got) != 2 || got[0] != "cmd/" || got[1] != "internal/" {
		t.Fatalf("paths.app = %v", got)
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
app_scan_roots = ["cmd"]
tools_scan_roots = ["tools"]
global_exclusions = [".git"]
generated_marker = "DO NOT EDIT."
generated_probe_limit = 128

[styleguide]
path = "STYLE.md"
requirement_id_format = "section_slug"

[imports]
local_prefix = "example.com/test"

[paths]
app = ["cmd/"]
application_port = ["internal/app/ports/"]
concrete_infra = ["internal/adapters/"]
domain = ["internal/domain/"]
domain_errors = ["internal/domain/errors.go"]
test_mocks = ["internal/testsupport/mocks/"]

[[file_sets]]
name = "markdown"
extensions = [".md"]
app_files = ["STYLE.md"]
tools_prefixes = ["tools/"]

[[language.backends]]
name = "go_app"
language = "go"
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
value = "--profile required"

[[control_plane.target_contracts]]
name = "lint"
recipe_line = "@$(STYLE_CMD) check $(LINT_FULL_ARGS)"

[[architecture.layers]]
name = "domain"
package_roots = ["internal/domain"]
may_import = ["domain"]

[[rules]]
rule_id = "repo/naming"
level = "required"
scope = "all"
requirement_ids = ["3.3.use-repository-name"]
`,
	)

	if _, err := Load(repoRoot); err == nil {
		t.Fatal("expected missing configured root marker to fail")
	}
}
