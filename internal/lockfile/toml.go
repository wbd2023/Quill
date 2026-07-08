package lockfile

import (
	"bytes"
	"fmt"
	"sort"

	codec "github.com/BurntSushi/toml"
)

type schemaConfig struct {
	SchemaVersion int             `toml:"schema_version"`
	Archives      []schemaArchive `toml:"archive"`
}

type schemaArchive struct {
	Tool    string            `toml:"tool"`
	Version string            `toml:"version"`
	Hashes  map[string]string `toml:"hashes"`
}

// Decode parses lockfile TOML source.
func Decode(source string) (lockfile Lockfile, err error) {
	var schema schemaConfig
	if _, err = codec.Decode(source, &schema); err != nil {
		return Lockfile{}, err
	}

	if schema.SchemaVersion != 1 {
		return Lockfile{}, fmt.Errorf(
			"unsupported lockfile schema_version %d (want 1)",
			schema.SchemaVersion,
		)
	}

	archives := make(map[string]Archive, len(schema.Archives))
	for _, entry := range schema.Archives {
		if _, duplicate := archives[entry.Tool]; duplicate {
			return Lockfile{}, fmt.Errorf("duplicate lockfile entry for %s", entry.Tool)
		}

		archives[entry.Tool] = Archive(entry)
	}

	return Lockfile{Archives: archives}, nil
}

// Encode writes lockfile content as TOML. Archive entries are sorted by tool ID for
// deterministic output (spurious git diffs otherwise, since map iteration is unordered).
func Encode(lockfile Lockfile) (contents string, err error) {
	schema := schemaConfig{
		SchemaVersion: 1,
	}

	tools := make([]string, 0, len(lockfile.Archives))
	for tool := range lockfile.Archives {
		tools = append(tools, tool)
	}
	sort.Strings(tools)

	for _, tool := range tools {
		schema.Archives = append(schema.Archives, schemaArchive(lockfile.Archives[tool]))
	}

	var buffer bytes.Buffer
	encoder := codec.NewEncoder(&buffer)
	encoder.Indent = ""
	if err = encoder.Encode(schema); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
