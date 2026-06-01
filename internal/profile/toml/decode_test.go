package toml_test

import (
	"testing"

	"ciphera/tools/internal/policy"
	"ciphera/tools/internal/profile/toml"
)

func TestDecodeReadsMultilineArrays(t *testing.T) {
	t.Parallel()

	config, err := toml.Decode(`schema_version = 1

[path_roles]
go_source = [
	"cmd/",
	"internal/",
]
`)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	requireEqual(t, "path_roles", policy.PathRoles{
		"go_source": {"cmd/", "internal/"},
	}, config.PathRoles)
}

func TestDecodeReadsTargetWorkingDirectory(t *testing.T) {
	t.Parallel()

	config, err := toml.Decode(`schema_version = 1

[targets.tools_go]
language = "go"
scope = "tools"
working_directory = "tools"
`)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	target, found := config.Targets.Lookup("tools_go")
	if !found {
		t.Fatal("expected tools_go target")
	}

	if target.WorkingDirectory != "tools" {
		t.Fatalf("working dir = %q", target.WorkingDirectory)
	}
}
