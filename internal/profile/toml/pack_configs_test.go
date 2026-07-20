package toml_test

import (
	"testing"

	"github.com/wbd2023/Quill/internal/profile/toml"
)

func TestDecodeReadsPackConfigs(t *testing.T) {
	t.Parallel()

	config, err := toml.Decode(`schema_version = 1

[packs]
enabled = ["go"]

[packs.go]
local_import_prefixes = ["ciphera"]

[packs.go.parameters]
secret_names = ["token"]
`)
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}

	pack, found := config.PackConfigs.Lookup("go")
	if !found {
		t.Fatal("expected go pack config")
	}

	parameters, ok := pack["parameters"].(map[string]any)
	if !ok {
		t.Fatalf("expected parameters map, got %T", pack["parameters"])
	}

	requireEqual(t, "pack parameters", map[string]any{
		"secret_names": []any{"token"},
	}, parameters)
}
