package profile

import "testing"

/* ------------------------------------------- Loading ------------------------------------------ */

func TestParseRejectsUnknownKeys(t *testing.T) {
	t.Parallel()

	_, err := parse(`profile_version = 1
unknown_key = true
`)
	if err == nil {
		t.Fatal("expected unknown key to fail")
	}
}

func TestParseDecodesMultilineArrays(t *testing.T) {
	t.Parallel()

	config, err := parse(`profile_version = 1

[paths]
go_source = [
	"cmd/",
	"internal/",
]
`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	got := config.Paths.LookupPatterns("go_source")
	if len(got) != 2 || got[0] != "cmd/" || got[1] != "internal/" {
		t.Fatalf("paths.go_source = %v", got)
	}
}

func TestParseDecodesSectionHeaderPolicy(t *testing.T) {
	t.Parallel()

	config, err := parse(`profile_version = 1

[formatting.section_headers]
required_min_lines = 100
short_file_max_lines = 79
overuse_threshold = 7
generic_names = ["Check", "Checks"]
structural_names = ["Types", "Helpers"]
`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}

	headers := config.Formatting.SectionHeaders
	if headers.RequiredMinLines != 100 ||
		headers.ShortFileMaxLines != 79 ||
		headers.OveruseThreshold != 7 ||
		len(headers.GenericNames) != 2 ||
		headers.GenericNames[0] != "Check" ||
		len(headers.StructuralNames) != 2 {
		t.Fatalf("unexpected section header policy: %#v", headers)
	}
}

func TestParseRejectsOldScanRootFields(t *testing.T) {
	t.Parallel()

	_, err := parse(`profile_version = 1

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
generated_probe_bytes = 128

[repository.scope_roots]
app = ["cmd"]
tools = ["tools"]
all = ["."]

[styleguide]
path = "STYLE.md"
requirement_id_scheme = "section_slug"

[formatting.section_headers]
required_min_lines = 100
short_file_max_lines = 79
overuse_threshold = 7
generic_names = ["Check", "Checks", "Misc", "Other"]
structural_names = ["Types", "Constants", "Helpers"]

[go]
local_import_prefixes = ["example.com/test"]

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
explicit_files = { app = ["STYLE.md"] }
path_prefixes = { tools = ["tools/"] }

[[language.backends]]
name = "application_go"
language = "go"
scope = "app"
workdir = "."
format_paths = ["cmd"]
check_paths = ["cmd"]

[vocabulary.go]
forbidden_type_suffixes = ["Repository"]
preferred_type_suffix = "Store"
forbidden_identifier_suffixes = ["Repository"]
preferred_identifier_suffix = "Store"

[vocabulary.shell]
forbidden_assignment_names = ["NC"]
preferred_assignment_name = "COLOUR_RESET"

[quality_surface]
driver = "make"

[quality_surface.make]
path = "mk/quality.mk"

[[quality_surface.make.required_variables]]
name = "LINT_REQUIRED_ARGS"
value = "--mode required"

[[quality_surface.make.required_targets]]
name = "lint"
recipe_line = "@$(STYLE_CMD) check $(LINT_FULL_ARGS)"

[[go.architecture.layers]]
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
