package profile

import (
	"bytes"

	"ciphera/tools/internal/policy"

	"github.com/BurntSushi/toml"
)

func Render(config policy.Config) (contents string, err error) {
	var buffer bytes.Buffer
	if err = toml.NewEncoder(&buffer).Encode(schemaFromPolicy(config)); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
