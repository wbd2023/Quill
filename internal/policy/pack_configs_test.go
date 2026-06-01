package policy_test

import (
	"testing"

	"ciphera/tools/internal/policy"
)

func TestPackConfigCloneCopiesArraysOfTables(t *testing.T) {
	t.Parallel()

	config := policy.PackConfig{
		"tables": []map[string]any{
			{"name": "first", "values": []string{"a"}},
		},
	}

	clone := config.Clone()
	config["tables"].([]map[string]any)[0]["name"] = "changed"
	config["tables"].([]map[string]any)[0]["values"].([]string)[0] = "b"

	cloneTable := clone["tables"].([]map[string]any)[0]
	requireEqual(t, "first", cloneTable["name"])
	requireEqual(t, "a", cloneTable["values"].([]string)[0])
}
