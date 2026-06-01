package toml

import (
	"fmt"
	"strings"

	"ciphera/tools/internal/policy"

	codec "github.com/BurntSushi/toml"
)

// Decode decodes style profile TOML source.
func Decode(source string) (config policy.Config, err error) {
	var schema schemaConfig
	metadata, err := codec.Decode(source, &schema)
	if err != nil {
		return policy.Config{}, err
	}

	for _, key := range metadata.Undecoded() {
		if strings.HasPrefix(key.String(), "packs.") {
			continue
		}

		return policy.Config{}, fmt.Errorf("unknown style.toml key %q", key.String())
	}

	return decodeConfig(schema)
}
