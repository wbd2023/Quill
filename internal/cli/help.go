package cli

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
)

/* ------------------------------------------ Root Help ----------------------------------------- */

// rootUsage renders the application usage, listing every visible command in declaration order.
// The column width is derived from the longest command name so usage stays aligned as commands are
// added or renamed.
func rootUsage(app *kong.Application) (usage string) {
	commands := visibleCommandNodes(app.Children)

	width := 0
	for _, command := range commands {
		if length := len(command.Name); length > width {
			width = length
		}
	}
	if width > 0 {
		width++
	}

	lines := []string{
		"usage:",
		"  quill <command> [flags]",
		"",
		"commands:",
	}
	for _, command := range commands {
		lines = append(lines, fmt.Sprintf("  %-*s %s", width, command.Name, command.Help))
	}

	lines = append(
		lines,
		"",
		"run `quill help <command>` or `quill <command> -h` to see command-specific flags",
		"",
	)
	return strings.Join(lines, "\n")
}

// rootUsageText renders root usage from a fresh command model.
func rootUsageText() (usage string) {
	parser, _, err := newParser()
	if err != nil {
		return ""
	}
	return rootUsage(parser.Model)
}

/* ---------------------------------------- Command Help ---------------------------------------- */

// commandUsage renders a single command's usage: invocation line, summary, then its flags sorted
// alphabetically for stable output regardless of field declaration order. Positional syntax is
// derived from Kong's model so help cannot drift from the grammar it describes.
func commandUsage(node *kong.Node) (usage string) {
	flags := visibleFlags(node.Flags)
	sort.Slice(flags, func(i int, j int) bool {
		return flags[i].Name < flags[j].Name
	})

	invocation := "quill " + node.Name
	for _, positional := range node.Positional {
		invocation += " " + positional.ShortSummary()
	}
	invocation += " [flags]"

	lines := []string{
		"usage:",
		"  " + invocation,
	}
	if node.Help != "" {
		lines = append(lines, "", node.Help)
	}
	if len(flags) > 0 {
		lines = append(lines, "", "flags:")
		for _, flag := range flags {
			lines = append(lines, formatFlagLines(flag)...)
		}
	}
	lines = append(lines, "")
	return strings.Join(lines, "\n")
}

/* ------------------------------------------- Helpers ------------------------------------------ */

// formatFlagLines renders a flag as two lines: the flag head on the first, its help (and default
// where set) indented on the second. Boolean flags omit the value placeholder and default.
func formatFlagLines(flag *kong.Flag) (lines []string) {
	head := "    --" + flag.Name
	if placeholder := flagPlaceholder(flag); placeholder != "" {
		head += " " + placeholder
	}

	help := flag.Help
	if def, ok := flagDefault(flag); ok {
		if help != "" {
			help += " "
		}
		help += fmt.Sprintf("(default %q)", def)
	}

	lines = []string{head}
	if help != "" {
		lines = append(lines, "        "+help)
	}
	return lines
}

// flagPlaceholder returns the value placeholder shown after a flag name, mirroring the type word
// convention (e.g. "string"). Boolean and counter flags take no placeholder.
func flagPlaceholder(flag *kong.Flag) (placeholder string) {
	if flag.IsBool() || flag.IsCounter() {
		return ""
	}
	switch flag.Target.Type().Kind() {
	case reflect.String:
		return "string"

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int"

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "uint"

	case reflect.Float32, reflect.Float64:
		return "float"

	default:
		return strings.ToLower(flag.Target.Type().String())
	}
}

// flagDefault reports the default value to display for a non-boolean flag. It reads the default
// literal declared on the tag (Value.Default), which is populated when the model is built and so is
// available even when help renders before Kong applies defaults. Empty and boolean defaults are
// suppressed.
func flagDefault(flag *kong.Flag) (value string, ok bool) {
	if flag.IsBool() || flag.IsCounter() || !flag.HasDefault {
		return "", false
	}

	if value = flag.Default; value == "" {
		return "", false
	}
	return value, true
}

func visibleCommandNodes(nodes []*kong.Node) (visible []*kong.Node) {
	visible = make([]*kong.Node, 0, len(nodes))
	for _, node := range nodes {
		if node.Hidden {
			continue
		}
		visible = append(visible, node)
	}
	return visible
}

func visibleFlags(flags []*kong.Flag) (visible []*kong.Flag) {
	visible = make([]*kong.Flag, 0, len(flags))
	for _, flag := range flags {
		if flag == nil || flag.Hidden || flag.Name == "help" {
			continue
		}
		visible = append(visible, flag)
	}
	return visible
}
