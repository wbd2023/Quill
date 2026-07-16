package toml_test

import (
	"testing"

	"ciphera/tools/internal/profile/toml"
)

func TestDecodeRejectsUnknownKeys(t *testing.T) {
	t.Parallel()

	_, err := toml.Decode(`schema_version = 1
unknown_key = true
`)
	if err == nil {
		t.Fatal("expected unknown key to fail")
	}
}

func TestDecodeRejectsOldKeys(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "profile version",
			input: "profile_version = 1\n",
		},
		{
			name: "rule level",
			input: `schema_version = 1

[[rules]]
id = "text/ascii"
level = "required"
scope = "all"
requirement_ids = ["0.1.example"]
`,
		},
		{
			name: "target workdir",
			input: `schema_version = 1

[targets.tools_go]
language = "go"
scope = "tools"
workdir = "tools"
`,
		},
		{
			name: "repository scan roots",
			input: `schema_version = 1

[repository]
app_scan_roots = ["cmd"]
`,
		},
		{
			name: "styleguide table",
			input: `schema_version = 1

[styleguide]
path = "STYLE.md"
`,
		},
		{
			name: "file set excluded names",
			input: `schema_version = 1

[file_sets.spelling]
excluded_names = ["LICENSE"]
`,
		},
		{
			name: "file set excluded name prefixes",
			input: `schema_version = 1

[file_sets.spelling]
excluded_name_prefixes = ["LICENSE."]
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, err := toml.Decode(test.input)
			if err == nil {
				t.Fatal("expected old key to fail")
			}
		})
	}
}
