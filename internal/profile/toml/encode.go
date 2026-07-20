package toml

import (
	"bytes"
	"strings"

	"github.com/wbd2023/Quill/internal/policy"

	codec "github.com/BurntSushi/toml"
)

// Encode encodes config as style profile TOML.
func Encode(config policy.Config) (contents string, err error) {
	var buffer bytes.Buffer
	encoder := codec.NewEncoder(&buffer)
	encoder.Indent = ""
	if err = encoder.Encode(encodeConfig(config)); err != nil {
		return "", err
	}

	return formatEncodedTables(buffer.String()), nil
}

func formatEncodedTables(contents string) (formatted string) {
	lines := strings.SplitAfter(contents, "\n")
	kept := make([]string, 0, len(lines))

	for index, line := range lines {
		table, ok := tableHeader(line)
		if ok && hasChildTable(table, lines[index+1:]) {
			continue
		}

		if ok &&
			len(kept) > 0 &&
			strings.TrimSpace(kept[len(kept)-1]) != "" {
			kept = append(kept, "\n")
		}

		kept = append(kept, line)
	}

	return strings.Join(kept, "")
}

func hasChildTable(parent string, lines []string) (found bool) {
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		child, ok := tableHeader(line)
		return ok && strings.HasPrefix(child, parent+".")
	}

	return false
}

func tableHeader(line string) (name string, found bool) {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "[[") ||
		!strings.HasPrefix(line, "[") ||
		!strings.HasSuffix(line, "]") {
		return "", false
	}

	return strings.TrimSpace(line[1 : len(line)-1]), true
}
