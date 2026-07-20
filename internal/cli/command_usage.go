package cli

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"strings"
)

func commandUsage(name string, summary string, flagSet *flag.FlagSet) (usage string) {
	var buffer bytes.Buffer

	lines := []string{
		"usage:",
		fmt.Sprintf("  quill %s [flags]", name),
	}
	if summary != "" {
		lines = append(lines, "", summary)
	}

	if flagSet != nil {
		flagText := strings.Trim(flagUsages(flagSet), "\n")
		if flagText != "" {
			lines = append(lines, "", "flags:", indentBlock(flagText, "  "))
		}
	}

	buffer.WriteString(strings.Join(lines, "\n"))
	buffer.WriteByte('\n')
	return buffer.String()
}

func indentBlock(value string, prefix string) (indented string) {
	lines := strings.Split(value, "\n")
	for index, line := range lines {
		lines[index] = prefix + line
	}

	return strings.Join(lines, "\n")
}

func flagUsages(flagSet *flag.FlagSet) (usage string) {
	var buffer bytes.Buffer

	flagSet.SetOutput(&buffer)
	flagSet.PrintDefaults()
	flagSet.SetOutput(io.Discard)
	return strings.ReplaceAll(buffer.String(), "\t", "  ")
}
