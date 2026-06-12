package report

import (
	"encoding/json"
	"io"
)

func writeJSON(writer io.Writer, value any) (err error) {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}
