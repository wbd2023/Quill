package lockfile

import (
	"bytes"
	"fmt"

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

// Encode writes lockfile content as TOML.
func Encode(lockfile Lockfile) (contents string, err error) {
	schema := schemaConfig{
		SchemaVersion: 1,
	}

	for _, archive := range lockfile.Archives {
		schema.Archives = append(schema.Archives, schemaArchive(archive))
	}

	var buffer bytes.Buffer
	encoder := codec.NewEncoder(&buffer)
	encoder.Indent = ""
	if err = encoder.Encode(schema); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
