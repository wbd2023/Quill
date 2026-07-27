package toml

import (
	"fmt"
	"strings"

	"github.com/wbd2023/quill/internal/policy"

	codec "github.com/BurntSushi/toml"
)

// Decode decodes style profile TOML source.
func Decode(source string) (config policy.Profile, err error) {
	var schema schemaConfig
	metadata, err := codec.Decode(source, &schema)
	if err != nil {
		return policy.Profile{}, err
	}

	for _, key := range metadata.Undecoded() {
		if strings.HasPrefix(key.String(), "packs.") {
			continue
		}

		return policy.Profile{}, fmt.Errorf("unknown quill.toml key %q", key.String())
	}

	return decodeConfig(schema)
}
