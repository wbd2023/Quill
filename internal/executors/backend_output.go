package executors

import "strings"

func appendExecutorOutput(builder *strings.Builder, output string) {
	output = strings.TrimSpace(output)
	if output == "" {
		return
	}

	if builder.Len() > 0 {
		builder.WriteString("\n")
	}

	builder.WriteString(output)
}
